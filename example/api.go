package example

// CreateVPC : ...
// @ExternalAPI ( server="vpc", url="/v2.0/{project_id}/vpc", method="POST", code="200,201")
func CreateVPC() error {
	return nil
}

// DeleteVPC : ...
// @ExternalAPI (server="vpc", url="/v2.0/{project_id}/vpc", method="DELETE", code="200,404", body="not found")
func DeleteVPC() error {
	return nil
}
