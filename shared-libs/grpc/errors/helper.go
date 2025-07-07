package grpc_errors

import "google.golang.org/grpc/status"

func ExtractErrorDetails(err error) map[string]interface{} {
	if st, ok := status.FromError(err); ok {
		details := make(map[string]interface{})
		details["code"] = st.Code().String()
		details["message"] = st.Message()

		// Extract proto details if any
		for _, detail := range st.Details() {
			// Handle your specific proto types here
			_ = detail
		}

		return details
	}
	return nil
}
