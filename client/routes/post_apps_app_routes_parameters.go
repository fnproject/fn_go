// Code generated by go-swagger; DO NOT EDIT.

package routes

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/fnproject/fn_go/models"
)

// NewPostAppsAppRoutesParams creates a new PostAppsAppRoutesParams object
// with the default values initialized.
func NewPostAppsAppRoutesParams() *PostAppsAppRoutesParams {
	var ()
	return &PostAppsAppRoutesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPostAppsAppRoutesParamsWithTimeout creates a new PostAppsAppRoutesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPostAppsAppRoutesParamsWithTimeout(timeout time.Duration) *PostAppsAppRoutesParams {
	var ()
	return &PostAppsAppRoutesParams{

		timeout: timeout,
	}
}

// NewPostAppsAppRoutesParamsWithContext creates a new PostAppsAppRoutesParams object
// with the default values initialized, and the ability to set a context for a request
func NewPostAppsAppRoutesParamsWithContext(ctx context.Context) *PostAppsAppRoutesParams {
	var ()
	return &PostAppsAppRoutesParams{

		Context: ctx,
	}
}

// NewPostAppsAppRoutesParamsWithHTTPClient creates a new PostAppsAppRoutesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPostAppsAppRoutesParamsWithHTTPClient(client *http.Client) *PostAppsAppRoutesParams {
	var ()
	return &PostAppsAppRoutesParams{
		HTTPClient: client,
	}
}

/*PostAppsAppRoutesParams contains all the parameters to send to the API endpoint
for the post apps app routes operation typically these are written to a http.Request
*/
type PostAppsAppRoutesParams struct {

	/*App
	  name of the app.

	*/
	App string
	/*Body
	  One route to post.

	*/
	Body *models.RouteWrapper

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the post apps app routes params
func (o *PostAppsAppRoutesParams) WithTimeout(timeout time.Duration) *PostAppsAppRoutesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post apps app routes params
func (o *PostAppsAppRoutesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post apps app routes params
func (o *PostAppsAppRoutesParams) WithContext(ctx context.Context) *PostAppsAppRoutesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post apps app routes params
func (o *PostAppsAppRoutesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post apps app routes params
func (o *PostAppsAppRoutesParams) WithHTTPClient(client *http.Client) *PostAppsAppRoutesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post apps app routes params
func (o *PostAppsAppRoutesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithApp adds the app to the post apps app routes params
func (o *PostAppsAppRoutesParams) WithApp(app string) *PostAppsAppRoutesParams {
	o.SetApp(app)
	return o
}

// SetApp adds the app to the post apps app routes params
func (o *PostAppsAppRoutesParams) SetApp(app string) {
	o.App = app
}

// WithBody adds the body to the post apps app routes params
func (o *PostAppsAppRoutesParams) WithBody(body *models.RouteWrapper) *PostAppsAppRoutesParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the post apps app routes params
func (o *PostAppsAppRoutesParams) SetBody(body *models.RouteWrapper) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *PostAppsAppRoutesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param app
	if err := r.SetPathParam("app", o.App); err != nil {
		return err
	}

	if o.Body == nil {
		o.Body = new(models.RouteWrapper)
	}

	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
