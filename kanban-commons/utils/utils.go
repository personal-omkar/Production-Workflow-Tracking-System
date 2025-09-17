package utils

import (
	"archive/tar"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"log/slog"
	"math/big"
	mran "math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"irpl.com/kanban-commons/model"
)

const (
	SECRET_KEY = "vsys_jwt_token"
)

const DefaultRestHost string = "0.0.0.0" // Default port if not set in env
const DefaultRestPort string = "4200"    // Default port if not set in env

func init() {
	RestHost = os.Getenv("RESTSRV_HOST")
	if strings.TrimSpace(RestHost) == "" {
		RestHost = DefaultRestHost
	}

	RestPort = os.Getenv("RESTSRV_PORT")
	if strings.TrimSpace(RestPort) == "" {
		RestPort = DefaultRestPort
	}

	RestURL = JoinStr("http://", RestHost, ":", RestPort)

	// map actual status to Status tobe displayed
	StatusMap["Dispatch"] = "Dispatched"
	StatusMap["dispatch"] = "Dispatched"
	StatusMap["Pending"] = "Pending"
	StatusMap["pending"] = "Pending"
	StatusMap["creating"] = "Created"
	StatusMap["Creating"] = "Created"
	StatusMap["Approved"] = "Approved"
	StatusMap["approved"] = "Approved"
	StatusMap["Reject"] = "Rejected"
	StatusMap["reject"] = "Rejected"
	StatusMap["InProductionProcess"] = "Production Line"
	StatusMap["InProductionLine"] = "Lined Up"
	StatusMap["true"] = "True"
	StatusMap["True"] = "True"
	StatusMap["false"] = "False"
	StatusMap["False"] = "False"
}

// concatenate all provided strings and return single resultant string
func JoinStr(sentences ...string) (result string) {
	var final strings.Builder
	for _, sentence := range sentences {
		final.WriteString(sentence)
	}
	result = final.String()
	return
}

// compares the data in both interface and returns true if both interface content is identical else false
func CompareData(first interface{}, second interface{}) (res bool) {
	var fir []byte
	var sec []byte

	fir, firOk := first.([]byte)
	if !firOk {
		byteArray, firErr := json.Marshal(first)
		if firErr != nil {
			res = false
		}
		fir = byteArray
	}

	sec, secOk := second.([]byte)
	if !secOk {
		byteArray, secErr := json.Marshal(second)
		if secErr != nil {
			res = false
		}
		sec = byteArray
	}
	res = bytes.Contains(fir, sec)

	return
}

// returns date corresponding to yyyy-mm-dd hh-mm-ss + nsec local(+530 IST)
func DateAndTime(year int, month time.Month, day int, hour int, min int, sec int, nsec int) time.Time {
	return time.Date(year, month, day, hour, min, sec, nsec, time.Local)
}

// Create Jwt token using mobile number and hard coded secret key
// Token expiry timing is kept as 1 hr
func CreateJwtToken(email string) (string, time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = email
	expire := time.Now().UTC().Add(time.Hour * 1)
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = time.Now().UTC().Unix()
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", time.Now(), err
	}
	return tokenString, expire, nil
}

// ValidateJwtToken checks the validity of the JWT token
func ValidateJwtToken(tokenString string) bool {
	secretKey := []byte(SECRET_KEY)

	// log.Println("ValidateJwtToken: tokenString: ", tokenString)

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method: %v", token.Header["alg"])
			return nil, nil // Should return an error here instead of nil, nil if the signing method is not as expected
		}
		return secretKey, nil
	})

	if err != nil {
		log.Println("Token parsing error, token- "+tokenString+" error- ", err)
		return false
	}

	// Validate token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if the token is expired
		if exp, ok := claims["exp"].(float64); ok {
			// Convert to int64, then compare with current time
			if int64(exp) > time.Now().UTC().Unix() {
				return true // Token is valid and not expired
			}
		}
	}

	return false // Token is either invalid or expired
}

func GenerateOtp(max int) (int, string, error) {
	var digitsTable = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	var b = make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		return 0, "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = digitsTable[int(b[i])%len(digitsTable)]
	}
	byteToInt, err := strconv.Atoi(string(b))
	if err != nil {
		return 0, "", err
	}
	// hard coding, need to replace this with near to optimal logic
	if max == 6 && byteToInt < 100000 {
		return GenerateOtp(max)
	}
	return byteToInt, string(b), nil
}

// Nop not operation used to avoid send html comments
func Nop(strs ...string) {
}

// Respond,ailure if http status code is not equal to 200 or 201, then return failure
func RespondFailure(resp *http.Response, w http.ResponseWriter) {
	// Check the status code of the response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		// log.Printf("Received non-success status code: %d", resp.StatusCode)

		// Read the response body from the unsuccessful request
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			http.Error(w, "Failed to get requested data", resp.StatusCode)
			return
		}

		// Write the received status code and response body back to the client
		w.WriteHeader(resp.StatusCode)
		w.Write(responseBody)
		return
	}
}

// IsAdministrator reads request headers and returns true if the user is an administrator.
func IsAdministrator(request *http.Request) bool {
	role := request.Header.Get("X-Custom-Role")
	if role == "" {
		log.Printf("IsAdministrator - role header is empty")
		return false
	}
	return role == "Administrator"
}

func Error(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func IF(cond bool, result ...string) (res string) {
	if cond {
		res = JoinStr(result...)
	}
	return
}

func IfElse(cond bool, trueVal, falseVal string) (res string) {
	if cond {
		res = trueVal
	} else {
		res = falseVal
	}
	return
}

// ParseInt tries to parse a string into an int, returns defaultValue if parsing fails
func ParseInt(value string, defaultValue int) (int, error) {
	if value == "" {
		return defaultValue, nil
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue, err
	}
	return parsedValue, nil
}

const any = "k|qY=*Es3ko}q$Vb7%+edD-g4MH(UQma{iC?c7;Y1fhIDV`o&i)~pYM;QN~]ZmsL_S/[!5nz9scSUhn=?}" +
	"{rx*G&#8j7H[ERnJg-dH&Uv7nfG{IF$*?}W.bOL`1PjkSvur@H@?oSaqQhkimrh|g`w1wOue?$f%Un/LPH7S(/y%9~+$7Ckl" +
	"kpI2pUg]Jp&W@zW#6*(c_`u^<r-xim#YWCnDd<oZfbgPkFGV<g2Du_)avx_D?-fk}+!,:|38(<O!xEIqc3Ld(xib$P|Y;V" +
	"PQFs]%j^}YM.=ja0*3R1O5rOZw<z7QJ/T9`JC~ygZ3.!:#q;1Lp0o(j{6)o0^I4vr-7,g;4h/-P:!ceX~qdy|k(mfMP4$H>N{." +
	"N%?V-PrhrX#}:gR~&32}7u12+]>0AuNP)!Y1Na!9ZoT9%[+~.eZ@(QPNq{q#ddep9:`]nLs]dj%D#>5Cs|)Z?oo4;kJnbMkc" +
	"sVD{,%~v+@UHvc403rl%mDZ#z{!Pk$&nAyj,N6)%9B`0~G}KvuyzM@]#ZjiS*L{g{H]O7D5>hmxq:#7OscF5J2RYq+*KbtPXhX!" +
	",]K5}}MMlE@y;,<Nqs@7]!Mz|W8%>z.<qsa)s.#B[[Dm.Ij#/@us%B5x-5A}m_91EqbZMtYiQ.Ll&>QK:9LV{ekgb6~o[K,o3K" +
	"l8zSK5^0kCc:UFAP!tE%hv7lE^BRu-UTV@Y`YDj)&bk/B}e/I7GXK>m9QEts$RL3/=$9d-RI+tNK:jhFYVoD>L6ZDW~y=*Hr,H" +
	"nR;st2jo-U<FSWB?gOj6;]Lh1Z3MN*ht}TEvLG:#jP7/@l[9p$:KGCvz*^<<?pnomp+BS$.#JtyQW&!a{mA%K>?-nH8],`)?ONn" +
	"#l+U+I0A3DdAbX7R[u/EIhp[=NGsItCn$xLL]H#{C*&EyTAQcRQQiKo{oy0-3ho.=+,-kEPn%yQeVLPwI+/&,Y.1dYIeW8xf~[r_" +
	"UKZI;u7nXqqbVzX9%O{ZMCU&f=Tz}}(%pk_Jtm6rqAv5slBL}bmXjXi7_`:-u%ka3iEN0qMw?n278v+6zM$(MS]EscpsN>Rl#zi.^" +
	"Lr,Bm1}RxMg<;5*kWfM~V}x}6?<_].}`8]$}d1[BavUk4xCtTppe9pM[PGNI7.5X|y0lXnVukX~0!$g#x7y+~`U8I>7*q!kj+" +
	"4N%Z6`XFgQYZUd)sqKun;@q,*D/5F26G/u#|UNQITAJ?|J;Xz{wT#%@wp#$;7B1%*p(N]U]Eph>:9o!!-t0G#ySSA8OCT{hiC46" +
	"!YoBWp>4o7?/4=]g8A{/*v3N)lg+Yw]/0GoR<h}BhD^,8An(8!;GQ3=G{Du=MQYz$Mm;!N_F;cZ?Ia6zZ^L>Pi5T<ODPfh1qch}" +
	"TPYj5`FG+R4{R,>*`Sj?EuZ2)E)$wSIKFGPOcX=K8D/Re4RIW14{)d,#*{@EN149P.4wcXkw`z?IZ:aXq}|K{^+.opU*xHN`@m;3" +
	"fG]-XD^Zyof-ZCS;8>zU@>M]aT030@qYKUJj#WYp.~1uH|F8dz&?WheiC[LFEXYIf3pxm7V;~KyL/+,[}8|/=Q8I|kU8!?n&o" +
	"aJo5)rHo=6f5yGq}[f+02_+D?`u)|mH]uc$WG*k/&4B3@PsV^D-B4oqVV(.lxQ.w6hX{Y@|{iO]64e%sZkVly#rbgs5ck#.K9b1" +
	"`B8a&bz_H-hrDkU4KP[}tvj,6>7oJcLgy+[bP:NT|dTO$abK:03QLnI]XLJbb!9eg:FYPb*]#9pCyd1ND_T}kLz[X~2guac3G.u" +
	"y1zbv/W.>DktT<,.BE4J@oOH|vX|q*CR{v-q1Al{m3/(K4mRx-(+$IilEPJ7+l@;TQ6P-:QB}!LB.Ub0?(6>X/J<|jY4GXKK?*7+p" +
	"Cioqy/ij+N.c-Y<UNHz9>&%DOnkvO).*%HPT&HC:`~eBRu+U+CiFNkDfrMLk{@aZi5ZhN@S{2nO,Fs$(kKwMk<zS/L~nOZwfaDro" +
	"j-F[#,ce5W(;Xfm#mh,{RlNiaIW7J7?D+qCSJ[8q6Nb#x*!rPkN+^|}YU1JLPooaj&ELQ9CxP),R(4IajNG]&VEl_;oo3g+_ygkfFA" +
	"iEu9jTD_E;hUM<#.!FaS5/6[c(]$].(Ys;a6+dni3B/iUcjf*Xz.[DdFf?*RmoOqT!P(1&@{g_*cw_GH>9MiYWGIw^;[GrIkl3<>Ry" +
	"G]ntTplXD}gsn@AQ&*T?)S>Tw2=E~nX)f-1f]fQ4tJyNDYjGVpr&qc&0!oruxvXHDQBNJq.)bO=]$hbNYq)FBwsh#(DOot2Y&!uJ1=" +
	"^/I7QNs:+IZF.c]>DN]5(`I4mc>!DX4V=rw<~Gp}|PjBfk3@N{g.0F{(L!p1`PW,@r0ANm<3pN^S#r>x<FFu{L%}G&Aw<=(;F=<4MJ" +
	")3[+pl7;asqe0g,!px~dGL+Al~7:Q:6Em5JRTq19LqNQ>Nt(R=N<4S@!n1BSU|sq+KJ?#8IO5(#5Y|PQOkxz{C~2lc=j|71v0W|=Tt" +
	"bOe0^|f*Li0xgWCKNnmWGtDL?5d-b}E&Z9KA~xTfFEk11O%vlvPu|?im:&^%HA&;(x!9CH$_2KuX@,IKG`.peD+1Ivx4`#,oPbI>?p" +
	"n[1ER.LeUvc(*n0|/UfqDKdJ,oUV5NNj7wQgZ,ENuR9q]?jBxvHH1m}K;O+55l;LW=Q[s^9xeX#FHH~Yp:sG5EMC~s_jH-F2(Y+P3v#" +
	"!BMj9qMkoNxH_[!=gKwz4ddY~e!N)_`cXh<Z>l#M+z;P$OCc}4:}ShyqA}/,dcp|uCZ5JRBK@{bQ+*lJM|ooG3TCE:OcVrjiIIR/fNV-" +
	"cmo{TfL6>v+t<Cb@o1<g3:CoNl>/~C&Al%b=;>E]+lZ^m@$F^hc)6R(/2s+rN*<Gata^ELT/YvDp;d1-]!|Y_J:,Ljpm`>KcK%lzy1LE" +
	"SH*jpnuM/<6;$1#^=k=/P%gDq,us5G~t&njw;f(_^WE=8$651TIw,eh4;}v~*f&yTsd38<@^Q.AxlA_Rm:;i9n.ST|gjvLAVsyD{ZpC" +
	"IaU-5t/$<?]fta&E&w(VEhyL7/K|wd0TG,o9K<`tf3|u[B,|0_-luQ:1hM6%Ut{p%?VruBs|AjzWOF>olt)&q7$M^7w#Kd{hy,|^~M*" +
	"4P)vYC|x8Y8^78NgizD>^V1#RrP8Y<}U2,F(Gnun,Z:d4`C0GtDFSM=GyG&u,lU)CHB]CR8<sn=1D[&H{PB2oP>b#" +
	"wfeWZ`r^$C3&]5hvB)[1Vpf*d0FaT}"

// start + length < 3200
func GetPredefKey(start, length int) string {
	if start+length > len(any) {
		start = 0
	}
	return string([]byte(any)[start : start+length])
}

// text := []byte("My Super Secret Code Stuff 64")
// key := []byte("passphrasewhichneedstobe64bytes!")
func KeylessEncrypt64(data []byte) []byte {

	s2 := mran.NewSource(time.Now().UnixNano())
	r1 := mran.New(s2)
	r2 := r1.Intn(32000)

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, uint16(r2))
	if err != nil {
		log.Println("binary.Write failed:", err)
	}
	key := []byte(GetPredefKey(r2, 64))

	encr := AESEncrypt(data, key)
	encr = append(buf.Bytes(), encr...)
	final := hex.EncodeToString(encr)

	return []byte(final)
}

func KeylessDecrypt64(data []byte) (final []byte, err error) {
	// Decode the hex string
	if final, err = hex.DecodeString(string(data)); err == nil && len(final) > 3 {
		u := binary.BigEndian.Uint16(final[:2])
		// Fetch a 64-byte key instead of 32-byte
		key := []byte(GetPredefKey(int(u), 64))
		final, err = AESDecrypt(final[2:], key)
	} else {
		if err != nil {
			log.Println("KeylessDecrypt64", err.Error())
		}
	}
	return
}

//text := []byte("My Super Secret Code Stuff")
//key := []byte("passphrasewhichneedstobe32bytes!")

func KeylessEncrypt(data []byte) []byte {

	s2 := mran.NewSource(time.Now().UnixNano())
	r1 := mran.New(s2)
	r2 := r1.Intn(32000)

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, uint16(r2))
	if err != nil {
		log.Println("binary.Write failed:", err)
	}
	key := []byte(GetPredefKey(r2, 32))

	encr := AESEncrypt(data, key)
	encr = append(buf.Bytes(), encr...)
	final := hex.EncodeToString(encr)

	return []byte(final)
}

func KeylessFileEncrypt(inputFilePath string, outputFilePath string) error {

	// Read data from input file
	data, err := os.ReadFile(inputFilePath)
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		return err
	}

	s2 := mran.NewSource(time.Now().UnixNano())
	r1 := mran.New(s2)
	r2 := r1.Intn(32000)

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, uint16(r2))
	if err != nil {
		log.Println("binary.Write failed:", err)
	}
	key := []byte(GetPredefKey(r2, 32))

	encr := AESEncrypt(data, key)
	encr = append(buf.Bytes(), encr...)
	final := hex.EncodeToString(encr)

	// Write the result to the output file
	err = os.WriteFile(outputFilePath, []byte(final), 0644)
	if err != nil {
		log.Printf("Failed to write file: %v", err)
		return err
	}

	return nil
}

func KeylessFileDecrypt(inputFilePath string, outputFilePath string) error {
	// Read file
	data, err := os.ReadFile(inputFilePath)
	if err != nil {
		return err
	}

	// Perform decryption
	var final []byte
	if final, err = hex.DecodeString(string(data)); err == nil && len(final) > 3 {
		u := binary.BigEndian.Uint16(final[:2])
		key := []byte(GetPredefKey(int(u), 32))
		final, err = AESDecrypt(final[2:], key)
		if err != nil {
			return err
		}
	} else {
		if err != nil {
			log.Println("KeylessDecrypt", err.Error())
		}
	}

	// Write decrypted data to output file
	err = os.WriteFile(outputFilePath, final, 0644)
	if err != nil {
		return err
	}

	return nil
}

func KeylessDecrypt(data []byte) (final []byte, err error) {

	if final, err = hex.DecodeString(string(data)); err == nil && len(final) > 3 {
		u := binary.BigEndian.Uint16(final[:2])
		key := []byte(GetPredefKey(int(u), 32))
		final, err = AESDecrypt(final[2:], key)
	} else {
		if err != nil {
			log.Println("KeylessDecrypt", err.Error())
		}
	}
	return
}

func AESEncrypt(data, key []byte) []byte {

	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher(key)
	// if there are any errors, handle them
	if err != nil {
		fmt.Println(err)
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		fmt.Println(err)
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	return gcm.Seal(nonce, nonce, data, nil)
}

func AESDecrypt(data, key []byte) ([]byte, error) {

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		fmt.Println(err)
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}
	return plaintext, err
}

func ProcessTarFile(oldfilepath string) error {

	// Create a temporary directory
	tmpDir, _ := os.MkdirTemp("", "example")

	//tmpDir, _ := os.CreateTemp("", "example")
	defer os.RemoveAll(tmpDir)

	// Open the tar file
	oldfile, err := os.Open(oldfilepath)
	if err != nil {
		return err
	}
	defer oldfile.Close()

	tr := tar.NewReader(oldfile)

	// Extract all the files into the temporary directory
	for {

		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		target := filepath.Join(tmpDir, hdr.Name)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			f.Close()
		default:
			return fmt.Errorf("unknown type: %c in %s", hdr.Typeflag, hdr.Name)
		}

	}

	// Iterate over the files in the temporary directory and encrypt any JSON files
	err = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Encrypt the file contents
			encryptedData := KeylessEncrypt(data)

			// Replace the original file with the encrypted data
			err = os.WriteFile(path, encryptedData, info.Mode())
			if err != nil {
				return err
			}
		}

		return nil

	})

	if err != nil {
		return err
	}

	// Create a new tar file from the temporary directory
	newfile, err := os.Create(strings.Replace(oldfilepath, ".tar", ".cry.tar", 1))
	if err != nil {
		return err
	}
	defer newfile.Close()

	tw := tar.NewWriter(newfile)
	defer tw.Close()

	err = filepath.Walk(tmpDir, func(file string, fi os.FileInfo, err error) error {

		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(file, tmpDir+"/")

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(tw, f)
		return err

	})

	return err
}

func ReverseTarFile(oldfilepath string) error {
	// Create a temporary directory
	tmpDir, _ := os.MkdirTemp("", "example")
	defer os.RemoveAll(tmpDir)

	// Open the tar file
	oldfile, err := os.Open(oldfilepath)
	if err != nil {
		return err
	}
	defer oldfile.Close()

	tr := tar.NewReader(oldfile)

	// Extract all the files into the temporary directory
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		target := filepath.Join(tmpDir, hdr.Name)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			f.Close()
		default:
			return fmt.Errorf("unknown type: %c in %s", hdr.Typeflag, hdr.Name)
		}
	}

	// Iterate over the files in the temporary directory and decrypt any JSON files
	err = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Decrypt the file contents
			decryptedData, _ := KeylessDecrypt(data)

			// Replace the encrypted file with the decrypted data
			err = os.WriteFile(path, decryptedData, info.Mode())
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Create a new tar file from the temporary directory
	newfile, err := os.Create(strings.Replace(oldfilepath, ".cry.tar", ".tar", 1))
	if err != nil {
		return err
	}
	defer newfile.Close()

	tw := tar.NewWriter(newfile)
	defer tw.Close()

	err = filepath.Walk(tmpDir, func(file string, fi os.FileInfo, err error) error {
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(file, tmpDir+"/")

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(tw, f)
		return err
	})

	return err
}

func ChangeDate(date string) {

	//cred := GetCredentials(credentials)
	dateCmd := `sudo /usr/bin/date -s "` + date + `"`

	//cmdStr := "echo '" + cred[1] + "' | sudo -S " + dateCmd
	cmd := exec.Command("/bin/sh", "-c", dateCmd)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}

	log.Printf("%s\n", out)

}

func GetCredentials(code string) (credentials []string) {

	cred, _ := KeylessDecrypt([]byte(code))
	credentials = strings.Split(string(cred), ":")
	return
}

func StartSettings(credentials string) {

	cred := GetCredentials(credentials)

	cmdStr := "echo '" + cred[1] + "' | sudo -S /usr/bin/gnome-control-center"
	cmd := exec.Command("/bin/sh", "-c", cmdStr)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}

	log.Printf("%s\n", out)

}
func EscapeHTMLUsingRegex(input string) string {
	s := html.EscapeString(input)
	return s
}

// GetCurrentYearName returns the current financial year name.
func GetCurrentYearName() string {
	currentYear := time.Now().Year()
	if time.Now().Month() >= 4 {
		return fmt.Sprintf("%d-%d", currentYear, currentYear+1)
	}
	return fmt.Sprintf("%d-%d", currentYear-1, currentYear)
}

// ConvertRecordsToJSON converts a slice of Records to JSON format
func ConvertRecordsToJSON(records []interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error converting records to JSON: %v", err)
	}
	return string(jsonData), nil
}

// Optional: Send JSON data to a REST API endpoint
func SendJSONToAPI(jsonData string, apiUrl string) error {
	resp, err := http.Post(apiUrl, "application/json", strings.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error sending JSON to API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	log.Println("Data sent successfully to API")
	return nil
}

// EmailValidator validates if a string is a valid email
func EmailValidator(email string) bool {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

// ExtractUsernameAndDomain extracts the username (part before '@') from an email
func ExtractUsernameAndDomain(email string) (string, string, error) {
	// Validate the email first
	if !EmailValidator(email) {
		return "", "", errors.New("invalid email format")
	}

	// Split the email into username and domain
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", "", errors.New("invalid email structure")
	}

	return parts[0], parts[1], nil
}

func SetVersion(Version, Build string) string {
	if len(Version) == 0 {
		Version = RandString(5)
		Build = "1"
	}
	return Version + "." + Build
}

// RandString ss
func RandString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return stringWithCharset(length, charset)
}

func stringWithCharset(length int, charset string) string {
	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		// Generate a random index using crypto/rand
		index, _ := rand.Int(rand.Reader, charsetLength)
		result[i] = charset[index.Int64()]
	}

	return string(result)
}

func FormatStringDate(input interface{}, formatType string) string {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05.000Z07:00",
		"2006-01-02T15:04:05.00Z07:00",
		"2006-01-02T15:04:05.0Z07:00",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05 -0700 MST",
	}
	var parsedDate time.Time
	var err error

	switch v := input.(type) {
	case string:
		for _, format := range formats {
			parsedDate, err = time.Parse(format, v)
			if err == nil {
				break
			}
		}
		if err != nil {
			return "Date or Time format not match"
		}
	case time.Time:
		parsedDate = v
	default:
		return "Unsupported input type"
	}

	// Return based on the formatType
	switch formatType {
	case "date":
		return parsedDate.Format("02.01.2006")
	case "time":
		return parsedDate.Format("15:04:05")
	case "date-time":
		return parsedDate.Format("02.01.2006  15:04:05")
	case "short-date-time":
		return parsedDate.Format("02.01.2006  15:04")
	default:
		return "Invalid format type"
	}
}

// Create SystemLog
func CreateSystemLogInternal(logData model.SystemLog) error {
	// Marshal the SystemLog struct into JSON
	jsonValue, err := json.Marshal(logData)
	if err != nil {
		slog.Error("Error marshaling SystemLog data", "error", err.Error())
		return err
	}

	// Make the POST request to the create-system-log API
	resp, err := http.Post(RestURL+"/create-system-log", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("Error making POST request to create SystemLog", "error", err.Error())
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body for create SystemLog", "error", err.Error())
		return err
	}

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		slog.Error("Failed to create SystemLog", "status", resp.StatusCode, "response", string(responseBody))
		return err
	}

	slog.Info("Successfully created SystemLog", "response", string(responseBody))
	return nil
}

// it takes output of GetRows(excelize package function) which is [][]string and cols->expectedLen
// return [][]string where each row has exactly expectedLen cols(cells)
func NormalizeRowsToLength(rows [][]string, expectedLen int) [][]string {
	var normalized [][]string

	for _, row := range rows {
		// Skip completely empty rows
		if len(row) == 0 {
			continue
		}

		// Pad row if it's shorter than expectedLen
		if len(row) < expectedLen {
			for i := len(row); i < expectedLen; i++ {
				row = append(row, "")
			}
		}

		normalized = append(normalized, row)
	}

	return normalized
}

// WaitForHTTPServer waits until server answer on "ipPort" (IP:PORT) using "/status" HTTP call
func WaitForHTTPServer(ipPort string) bool {
	waitTime := 0
	for {
		if IsServerStatusOK(ipPort) {
			return true
		}
		time.Sleep(2 * time.Second)
		waitTime += 2
		if waitTime%10 == 0 {
			log.Println("Waiting on", ipPort, "for", strconv.Itoa(waitTime), "seconds")
		}
	}
}

// IsServerStatusOK verifies if server answer "ok" in rest function "status"
func IsServerStatusOK(server string) (res bool) {

	if !strings.HasPrefix(server, "http://") && !strings.HasPrefix(server, "https://") {
		server = "http://" + server
	}
	return ServerStatus(server)
}

// ServerStatus verifies if server answer "ok" in rest function "status"
func ServerStatus(server string) (res bool) {
	return ServerPage(server + "/status")
}

// ServerStatus verifies if server answer "ok" in rest function "status"
func ServerPage(server string) (res bool) {

	waitTime := 0

	if !strings.HasPrefix(server, "http://") && !strings.HasPrefix(server, "https://") {
		server = "http://" + server
	}

	res = false
	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	if resp, err := client.Get(server); err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {

			res = true

		}

	} else {
		time.Sleep(2 * time.Second)
		// log.Println(server, "error:", err.Error())
		waitTime += 2
		if waitTime%10 == 0 {
			log.Println("Waiting on", server, "for", strconv.Itoa(waitTime), "seconds. Error:", err.Error())
		}
	}
	return
}

// getFirstFileInDirectory checks if the directory contains any files and returns the first file's name.

func GetFirstFileInDirectory(dirPath string) (bool, string, error) {

	entries, err := os.ReadDir(dirPath)

	if err != nil {
		log.Printf("error %v", err)
		return false, "", err

	}

	if len(entries) == 0 {

		return false, "", nil // Directory is empty

	}

	return true, entries[0].Name(), nil

}

// ImageFetcher struct to manage image paths dynamically.

type ImageFetcher struct {
	DirPath string

	DefaultImg string
}

// NewImageFetcher initializes an ImageFetcher instance.

func NewImageFetcher(dirPath, defaultImg string) *ImageFetcher {

	return &ImageFetcher{

		DirPath: dirPath,

		DefaultImg: defaultImg,
	}

}

// GetImagePath returns the correct file path if a file exists, otherwise returns a default image.

func (i *ImageFetcher) GetImagePath() string {

	if fileExist, firstFile, _ := GetFirstFileInDirectory(i.DirPath); fileExist {

		return i.DirPath + firstFile

	}

	return i.DefaultImg

}

func CopyFile(srcPath, destPath string) error {

	// Open the source file

	srcFile, err := os.Open(srcPath)

	if err != nil {

		slog.Error("Failed to open source file", "path", srcPath, "error", err)

		return err

	}

	defer srcFile.Close()

	// Create the destination file

	destFile, err := os.Create(destPath)

	if err != nil {

		slog.Error("Failed to create destination file", "path", destPath, "error", err)

		return err

	}

	defer destFile.Close()

	// Copy the file contents

	_, err = io.Copy(destFile, srcFile)

	if err != nil {

		slog.Error("Failed to copy file", "src", srcPath, "dest", destPath, "error", err)

		return err

	}

	// Flush to ensure all writes are completed

	err = destFile.Sync()

	if err != nil {

		slog.Error("Failed to sync destination file", "path", destPath, "error", err)

		return err

	}

	slog.Info("File copied successfully", "src", srcPath, "dest", destPath)

	return nil

}

// ClearDirectory removes all files from the given directory.

func ClearDirectory(dirPath string) error {

	files, err := os.ReadDir(dirPath)

	if err != nil {

		slog.Error("Failed to read directory", "path", dirPath, "error", err)

		return err

	}

	for _, file := range files {

		filePath := filepath.Join(dirPath, file.Name())

		err := os.RemoveAll(filePath)

		if err != nil {

			slog.Error("Failed to delete file", "file", filePath, "error", err)

			return err

		}

		slog.Info("Deleted file", "file", filePath)

	}

	return nil

}
