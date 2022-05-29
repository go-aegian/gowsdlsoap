package proxy

// WithMIMEMultipartAttachments is an Option to set SOAP MIME Multipart attachment support.
//
// Use Client.AddMIMEMultipartAttachment to add attachments of type MIMEMultipartAttachment to your SOAP request.
func WithMIMEMultipartAttachments() Option {
	return func(o *Options) {
		o.Mma = true
	}
}
