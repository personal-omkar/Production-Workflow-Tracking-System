package services

import (
	"strings"

	"irpl.com/kanban-commons/utils"
)

type Model struct {
	ID       string
	Type     string
	Title    string
	Sections []FormSection
	Buttons  string
	Alert    bool
	CSS      string
	JS       string
}

func (mf *Model) Build() string {

	model := utils.JoinStr(`
	
	<div class="container">
	`, mf.CSS, `
		<div class="modal fade" id="`, mf.ID, `" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1" aria-labelledby="`, mf.ID, `-label" role="dialog" aria-modal="true">
			<div class="modal-dialog modal-dialog-centered `, mf.Type, `" >
				<div class="modal-content">
					<div class="modal-header">
						<span style="color:#871a83;"><b>`, mf.Title, `</b></span>
						<button type="button" class="btn-close" data-bs-dismiss="modal"  aria-label="Close" ></button>
					</div>
					<div class="modal-body" data-group="`, strings.ToLower(strings.ReplaceAll(mf.ID, " ", "")), `" >
						`)

	for _, section := range mf.Sections {
		model = utils.JoinStr(model, section.GenerateSection())
	}

	model = utils.JoinStr(model, `
					</div>
					<div class="modal-footer text-muted d-flex justify-content-end  align-items-center">
						`, mf.Buttons, `
					</div>
				</div>
			</div>
		</div>	
	`, mf.JS, `
	</div>
	`)
	return model
}
