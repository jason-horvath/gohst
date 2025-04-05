package render

import "gohst/internal/auth"

type TemplateData struct {
    CSRFToken    	string      	// CSRF token for form protection
    AuthUser	  	*auth.User      // Pointer to the authenticated user (if any)
    FlashMessages 	map[string]any  // Slice for any flash messages (success/error)
	OldFormData    	map[string]any 	// Map for old input values (for form repopulation)
    Data         	any  			// Additional dynamic data specific to each page
}
