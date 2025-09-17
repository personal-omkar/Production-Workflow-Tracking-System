package report

import (
	"irpl.com/kanban-commons/utils"
)

type HtmlCardDialog struct {
	Data       string
	RecordId   string
	PageSource string
}

func (p *HtmlCardDialog) Build() string {

	//<!--html-->
	ret := utils.JoinStr(`
	<style>
		/* Basic styling for the content */
		body {
			font-family: Arial, sans-serif;
		}
		.table {
			width: 100%;
			border-collapse: collapse;
		}
		.table th, .table td {
			border: 1px solid #ccc;
			padding: 8px;
			text-align: left;
		}

		/* Modal styling */
		#report-dialog {
			display: none;
			position: fixed;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			background-color: rgba(0, 0, 0, 0.7);
			justify-content: center;
			align-items: center;
		}
		.modal-content {
			background-color: #fff;
			padding: 20px;
			width: 80%;
			max-width: 600px;
			border-radius: 8px;
			text-align: left;
			box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
		}
		.modal-buttons {
			text-align: right;
			margin-top: 20px;
		}
		.modal-buttons button {
			padding: 8px 16px;
			margin: 0 5px;
			border: none;
			cursor: pointer;
			border-radius: 4px;
		}
		.btn-print {
			background-color: #4CAF50;
			color: white;
		}
		.btn-cancel {
			background-color: #f44336;
			color: white;
		}

	    /*width: 99mm  = 374px*/
		/*height: 69mm = 260px*/
		.printTable tr td {
			border: 1px solid #000;
			border-spacing: 0;
			border-collapse: collapse;
			font-family: Calibri;
			font-size: 12px;
			padding: 0 5px;
		}

		p {
			margin: 0;
		}

	</style>

	<div id="report-dialog" class="modal fade" tabindex="-1" role="dialog" data-backdrop="static"  data-bs-backdrop="static"  data-bs-backdrop="static"  >
	<div class="modal-dialog modal-xl" role="document">
		<div class="modal-content">
			<div class="modal-body">
				<div class="container-fluid">
					<div class="row">
						<div class="col">
							<div id="printableContent">
							<style>
								/*width: 99mm  = 374px*/
								/*height: 69mm = 260px*/
								.printTable tr td {
									border: 1px solid #000;
									border-spacing: 0;
									border-collapse: collapse;
									font-family: Calibri;
									font-size: 12px;
									padding: 0 5px;
								}

								p {
									margin: 0;
								}
							</style>
								`, p.PageSource, `
							</div>
						</div>
					</div>
				</div>
			</div>

			<div class="modal-footer">
				<!-- <button class="btn-print" onclick="printContent()">Print</button> -->
				<button type="button" class="btn btn-primary btn-lg" style="background-color:#CF7AC2; border:none;" onclick="printContent()">Print</button>
				<!-- <button class="btn-cancel" onclick="closePrintDialog()">Cancel</button> -->
				<button type="button" class="btn btn-danger btn-lg" data-bs-dismiss="modal"  data-dismiss="modal">Cancel</button>
			</div>
		</div>
	</div>
	<script>
		// Function to open the print dialog
		function openPrintDialog() {
			document.getElementById('report-dialog').style.display = 'flex';
		}

		// Function to close the print dialog
		function closePrintDialog() {
			document.getElementById('report-dialog').style.display = 'none';
		}


		// Function to print the content
		function printContent() {
			const printContent = document.getElementById('printableContent').innerHTML;
			const printWindow = window.open('', '', 'width=800,height=600');

			`)
	//<!--!html-->
	ret += utils.JoinStr("printWindow.document.write(`")

	//<!--html-->
	ret += utils.JoinStr(`
				<html>
				<head>
					<title>Print Preview</title>
					<style>
					body { font-family: Arial, sans-serif; }
					.table { width: 100%; border-collapse: collapse; }
					.table th, .table td { border: 1px solid #ccc; padding: 8px; text-align: left; }
					</style>
				</head>
				<body onload="window.print(); window.close();">
					${printContent}
				</body>
				</html>
			`)
	//<!--!html-->
	ret += utils.JoinStr("`);")

	//<!--html-->
	ret += utils.JoinStr(`

			printWindow.document.close();
		}
	</script>
	</div>
    `)
	//<!--!html-->
	return ret
}
