package dto

type CreateRecruiterRequest struct {
	CompanyName 	string `json:"company_name"`
	CompanyWebsite  string `json:"company_website"`
	Message 		string `json:"message"`
}

type RecruiterResponse struct {
	ID 			int64 `json:"id"`
	Status 		string `json:"status"`
	Message 	string `json:"message"`
}

type ListRecruiterRequestsResponse struct {
	Requests 	[]RecruiterRequestSummary `json:"requests"`
	Total   	int64                     `json:"total"`
	Limit    	int                       `json:"limit"`
	Offset   	int                       `json:"offset"`
}

type RecruiterRequestSummary struct {
	ID 				int64 `json:"id"`
	UserID 			int64 `json:"user_id"`
	CompanyName 	string `json:"company_name"`
	CompanyWebsite  string `json:"company_website"`
	Message 		string `json:"message"`
	Status 			string `json:"status"`
}

type UpdateRecruiterRequestStatusRequest struct {
	ID int64 `json:"id"`
	Status string `json:"status"` // "approved" or "rejected"
}

