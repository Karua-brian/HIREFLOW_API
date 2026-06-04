package validator

import (
	"errors"
	"strconv"
)

// ValidateRegister performs basic validation on the registration input fields.
// It checks for required fields and validates the role.
func ValidateRegister(email, password string) error {
	
	if email == "" {
		return errors.New("email is required")
	}

	if password == "" {
		return errors.New("password is required")
	}

	/* if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	} */

	return nil
}

// ValidateLogin performs basic validation on the login input fields.
func ValidateLogin(email, password string) error {
	if email == "" {
		return errors.New("email is required")
	}

	if password == "" {
		return errors.New("password is required")
	}

	return nil
}

func ValidateRefreshToken(refreshToken string) error {
	if refreshToken == "" {
		return errors.New("refresh token is required")
	}

	return nil
}

func ValidateCreateJob(title, company, location, salary string) error {
	if title == "" {
		return errors.New("title is required")
	}

	if company == "" {
		return errors.New("company is required")
	}

	if location == "" {
		location = "Remote" // Default to "Remote" if location is not provided
	}

	if salary == "" {
		salary = "Not specified" // Default to "Not specified" if salary is not provided
	}

	

	return nil
}

func ValidateListJobs(limit, offset int) error {
	if limit < 0 {
		return errors.New("limit must be a non-negative integer")
	}

	if offset < 0 {
		return errors.New("offset must be a non-negative integer")
	}

	return nil
}

func ValidateApplyJob(jobID int64) error {
	if jobID <= 0 {
		return errors.New("invalid job ID")
	}

	return nil
}

func ValidateJWTHeader(authHeader string) error {
	if authHeader == "" {
		return errors.New("Authorization header is missing")
	}

	return nil
}

func ValidateRecruiterRequest(companyName, companyWebsite, message string) error {
	if companyName == "" {
		return errors.New("company name is required")
	}

	if companyWebsite == "" {
		return errors.New("company website is required")
	}

	if message == "" {
		return errors.New("message is required")
	}

	return nil
}

func ParsePaginationParams(limitStr, offsetStr string) (int, int, error) {
	limit := 10
	offset := 0

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit < 0 {
			return 0, 0, errors.New("invalid limit parameter")
		}
		limit = parsedLimit
	}

	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			return 0, 0, errors.New("invalid offset parameter")
		}
		offset = parsedOffset
	}

	return limit, offset, nil
}