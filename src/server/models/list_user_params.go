package serverModels

// ListUserParams is the params to use in a GetAllUsers request
type ListUserParams struct {
	AdminStatus string `query:"admin_status" url:"admin_status" enum:"admin,non_admin,all"        default:"all"`
	Privacy     string `query:"privacy"      url:"privacy"      enum:"private,friends,public,all" default:"all"`
	Membership  string `query:"membership"   url:"membership"   enum:"free,paid,all"              default:"all"`
}
