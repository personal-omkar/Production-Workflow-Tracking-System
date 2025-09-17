package basepage

import (
	"strings"

	"irpl.com/kanban-commons/utils"
	u "irpl.com/kanban-commons/utils"
)

// BasePage represents the entire basepage of the HTML.
type BasePage struct {
	ExtraHeaders string
	SideNavBar   SideNav
	TopNavBar    TopNav
	BgColor      string
	Username     string
	Content      string
	UserType     string
}

// SideNavItem represents an individual navigation item in the side menu.
type SideNavItem struct {
	ID        string
	Name      string
	Icon      string
	Link      string
	HideMenu  bool
	Items     []SideNavSubItem
	Style     string
	CustomOpt string
	Selected  bool
	UserType  UserType
	Hidden    bool
}

// TopNavItem represents an individual navigation item in the top menu.
type TopNavItem struct {
	ID        string
	Name      string
	Title     string
	Type      string
	Icon      string
	Link      string
	HideMenu  bool
	Width     string
	Style     string
	CustomOpt string
}

// MenuLink represents an subitem for Side Nav Item.
type SideNavSubItem struct {
	Name     string
	URL      string
	New      bool
	Selected bool
	Style    string
}

// SideNav represents the entire side navigation bar.
type SideNav struct {
	MenuItems []SideNavItem
}

// TopNav represents the entire side navigation bar.
type TopNav struct {
	UserType   string
	VendorName string
	MenuItems  []TopNavItem
}

type UserType struct {
	Admin    bool
	Operator bool
	Customer bool
}

// BuildMenu generates the HTML for the side menu.
func (m *SideNav) BuildMenu(userType string) string {
	var html strings.Builder
	html.WriteString(u.JoinStr(
		`<nav class="navbar navbar-light navbar-vertical navbar-expand-xl">
			<script>
				var navbarStyle = localStorage.getItem("navbarStyle");
				if (navbarStyle && navbarStyle !== 'transparent') {
				document.querySelector('.navbar-vertical').classList.add('navbar-${navbarStyle}');
				}
			</script>
	
			<a class="navbar-brand" href="">
				<div class="d-flex align-items-center py-3"><img class="me-2" src="` + u.DefaultsMap["home_logo_url"] + `" alt="" width="200" style="margin-left:1rem;" />
				</div>
			</a>

			<div class="collapse navbar-collapse" id="navbarVerticalCollapse">
				<div class="navbar-vertical-content scrollbar d-flex">
					<ul class="navbar-nav flex-column mb-3" id="navbarVerticalNav">`))

	for _, item := range m.MenuItems {

		m.AddMenuOption(item, &html, userType)

	}

	html.WriteString(`
					</ul>
					<div class="flex-grow-1 align-content-end">
						<p class="mb-5">v` + utils.Version + `</p>
					</div>
				</div>
			</div>
		</nav>
	`)
	return html.String()
}

// AddMenuOption generates the menu options for the side Menu Bar.
func (m *SideNav) AddMenuOption(item SideNavItem, html *strings.Builder, usertype string) {

	var backgorundColor string
	var foregorundColor string
	var hide string
	if item.Selected {
		backgorundColor = "background:#871a83"
		foregorundColor = "color:#ffffff"
	}
	if item.Hidden {
		hide = "d-none"
	} else {
		hide = ""
	}
	if len(item.Items) > 0 {

		html.WriteString(u.JoinStr(`
				<li class="nav-item rounded `, hide, `" style="`, backgorundColor, `">
				   <a class="nav-link dropdown-indicator" href="#`, item.Name, `" role="button" data-bs-toggle="collapse" aria-expanded="false" aria-controls="`, item.Name, `">
					   <div class="d-flex align-items-center ms-1">
					   <span class="nav-link-icon" style="`, item.Style, foregorundColor, `">
							   <span class="`, item.Icon, `"></span>
					   </span>
					   <span class="nav-link-text ps-1" style="`, item.Style, `">`, item.Name, `</span>
					   </div>
				   </a>
				`))
		html.WriteString(u.JoinStr(`<ul class="nav collapse" id="`, item.Name, `">`))
		for _, subItem := range item.Items {
			html.WriteString(u.JoinStr(`
					<li class="nav-item">
						<a class="nav-link" href="`, subItem.URL, `">
							<div class="d-flex align-items-center ms-1">
								<span class="nav-link-text ps-1" style="`, item.Style, foregorundColor, `">`, subItem.Name, `</span>
							</div>
						</a>
					</li>
				`))
		}
		html.WriteString(u.JoinStr(`</ul></li>`))
	} else {
		html.WriteString(u.JoinStr(`
			<li class="nav-item rounded `, hide, `" style="`, backgorundColor, `">
			   <a class="nav-link" href="`, item.Link, `" role="button">
				   <div class="d-flex align-items-center ms-1">
				   <span class="nav-link-icon" style="`, item.Style, foregorundColor, `">
						   <span class="`, item.Icon, `"></span>
				   </span>
				   <span class="nav-link-text ps-1" style="`, item.Style, foregorundColor, `">`, item.Name, `</span>
				   </div>
			   </a>
			</li>`))
	}

}

// Set navbar permission
func CheckdisabledNavItems(navitem []SideNavItem, value, separator string) []SideNavItem {
	// Split the value string into an array using the separator
	denyurlList := strings.Split(value, separator)

	// Update the Hidden field based on presence in the allowedValues array
	for i := range navitem {
		navitem[i].Hidden = false // Default to false
		for _, deny := range denyurlList {
			if navitem[i].Link == deny {
				navitem[i].Hidden = true
				break
			}
		}
	}

	// Return the updated menu links
	return navitem
}

// BuildMenu generates the HTML for the top menu.
func (m *TopNav) BuildMenu() string {
	style := ""
	var vendornamestyle string
	if m.VendorName != "" {
		vendornamestyle = ""
	} else {
		vendornamestyle = "p-3 d-none"
		style = "p-3"
	}
	var html strings.Builder
	html.WriteString(u.JoinStr(`
		<div class="col-7 	`, style, `">
			<div class="input-group mb-3">
				<div class="col-5">
					<button type="button" title="Vendor Name" class="btn w-100 bg-light bg-gradient `, vendornamestyle, `">
						`, m.VendorName, `
					</button>
				</div>
			</div>
		</div>
		`))

	for _, item := range m.MenuItems {
		m.AddMenuOption(item, &html)
	}

	return html.String()
}

// AddMenuOption generates the menu options for the top Menu Bar.
func (m *TopNav) AddMenuOption(item TopNavItem, html *strings.Builder) {

	if item.HideMenu {
		html.WriteString(``)
	}

	html.WriteString(`<div class="` + item.Width + `">`)

	if item.Type == "link" {
		if item.ID == "settings" {
			if m.UserType == "Admin" {
				html.WriteString(`<a class="btn w-100 ` + item.Style + `" href="` + item.Link + `" title="` + item.Title + `" ` + item.CustomOpt + `>
				<span class="` + item.Icon + `"></span>
			</a>`)
			} else {
				html.WriteString(`<a class="btn w-100 ` + item.Style + `"  title="` + item.Title + `" ` + item.CustomOpt + `>
				<span class="` + item.Icon + `"></span>
			</a>`)
			}

		} else {
			html.WriteString(`<a class="btn w-100 ` + item.Style + `" href="` + item.Link + `" title="` + item.Title + `" ` + item.CustomOpt + `>
			<span class="` + item.Icon + `"></span>
		</a>`)
		}

	} else if item.Type == "button" {
		html.WriteString(`<button type="button" class="btn w-100 ` + item.Style + `"  title="` + item.Title + `" ` + item.CustomOpt + `>
			` + item.Name + `
		</button>`)
	}

	html.WriteString(`</div>`)

}

// AddScriptCode injects <script></script> javascript tag into the <head> section of webpage
func (b *BasePage) AddScriptCode(script string) {
	b.ExtraHeaders += u.JoinStr(`<script  type="text/javascript">`, script, `</script>`)
}

// AddStyleCode injects <style></style> CSS tag into the <head> section of webpage
func (b *BasePage) AddStyleCode(style string) {
	b.ExtraHeaders += u.JoinStr(`<style>`, style, `</style>`)
}

// AddStyleCode injects <style></style> CSS tag into the <head> section of webpage
func (b *BasePage) AddStyleLink(style string) {
	b.ExtraHeaders += u.JoinStr(`<link href="` + style + `" rel="stylesheet" type="text/css">`)
}

// AddScriptLink injects <script src=”></script> CSS tag into the <head> section of webpage
func (b *BasePage) AddScriptLink(link string) {
	b.ExtraHeaders += u.JoinStr(`<script src="`, link, `"></script>`)
}

func (b *BasePage) Build() string {

	//<!--html-->
	basepage := u.JoinStr(`
<!DOCTYPE html>
<html lang="en-US" dir="ltr">

  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">


    <!-- ===============================================-->
    <!--    Document Title-->
    <!-- ===============================================-->
    <title>`+utils.DefaultsMap["name"]+`</title>

    <meta name="theme-color" content="#ffffff">
    <script src="/static/assets/js/config.js"></script>



    <!-- ===============================================-->
    <!--    Stylesheets-->
    <!-- ===============================================-->
    <link href="../static/assets/css/theme-rtl.min.css" rel="stylesheet" id="style-rtl">
    <link href="../static/assets/css/theme.min.css" rel="stylesheet" id="style-default">
	<link href="../static/assets/css/style.css" rel="stylesheet">
	<script src="../static/assets/js/index.js"></script>
	<link href="../static/vendors/flatpickr/flatpickr.min.css" rel="stylesheet" />
	<link href="/static/vendors/select2/select2.min.css" rel="stylesheet" />
	<link href="/static/vendors/select2-bootstrap-5-theme/select2-bootstrap-5-theme.min.css" rel="stylesheet" />


    <!-- ===============================================-->
    <!--    Javascript-->
    <!-- ===============================================-->
	<script src="/static/vendors/jquery/jquery.min.js"></script>
	<script src="/static/vendors/echarts/echarts.min.js"></script>
	<script src="/static/assets/js/echarts-example.js"></script>
	<script src="/static/vendors/bootstrap/bootstrap.min.js"></script>
	
	<style>
		.modal-size {
			max-width:1550px !important;
		}
		.pagination .page-item.active .page-link {
			background-color: #871A83;
			border: none;
			outline: none;
			color: white;
		}

		.pagination .page-item:not(.active) .page-link:hover {
			background-color: #ab71a2;
			border: none;
			outline: none;
			color: white;
		}
	</style>

    <script>
      var isRTL = JSON.parse(localStorage.getItem('isRTL'));
      if (isRTL) {
        var linkDefault = document.getElementById('style-default');
        linkDefault.setAttribute('disabled', true);
        document.querySelector('html').setAttribute('dir', 'rtl');
      } else {
        var linkRTL = document.getElementById('style-rtl');
        linkRTL.setAttribute('disabled', true);
      }
    </script>


	`, b.ExtraHeaders, `
  </head>


  <body class="overflow-hidden">

    <!-- ===============================================-->
    <!--    Main Content-->
    <!-- ===============================================-->
    <main class="main" id="top">
      <div class="container-fluid" data-layout="container-fluid">
        <script>
          var isFluid = JSON.parse(localStorage.getItem('isFluid'));
          if (isFluid) {
            var container = document.querySelector('[data-layout]');
            container.classList.remove('container');
            container.classList.add('container-fluid');
          }
        </script>
        `, b.SideNavBar.BuildMenu(b.UserType), `
		
		<!-- <div style="margin-left:15rem;">

		</div> -->
        <div class="content p-0">
		
			<div class="container">
				<div class="d-flex justify-content-start mt-3">
					<div class="row" style="width: 600vh;">
					`, b.TopNavBar.BuildMenu(), `
					</div>
				</div>
			</div>

			<!-- Use following div for all your content -->
			<div class="main-container">
				`, b.Content, `
		
			</div>
			<div id="dialog"> </div>
			<div id="additional-content">
			</div>
        </div>
      </div>
    </main>
    <!-- ===============================================-->
    <!--    End of Main Content-->
    <!-- ===============================================-->


	<!-- ===============================================-->
	<!--    JavaScripts-->
	<!-- ===============================================-->
	<script>
		$(document).on("keypress", ".modal form input", function(event) {
			if (event.key === "Enter" && $(this).val().trim() !== "") {
				event.preventDefault(); 
			}
		});
	</script>
	<script src="/static/vendors/popper/popper.min.js"></script>
	<script src="/static/vendors/anchorjs/anchor.min.js"></script>
	<script src="/static/vendors/is/is.min.js"></script>
	<script src="/static/vendors/typed.js/typed.js"></script>
	<script src="/static/vendors/select2/select2.min.js"></script>
	<script src="/static/vendors/select2/select2.full.min.js"></script>
	<script src="/static/vendors/fontawesome/all.min.js"></script>
	<script src="/static/assets/js/flatpickr.js"></script>
	<script src="/static/vendors/lodash/lodash.min.js"></script>
	<script src="/static/vendors/list.js/list.min.js"></script>
	<script src="/static/assets/js/theme.js"></script>

  </body>

</html>
	`)
	//	<!--!html-->
	return basepage

}
