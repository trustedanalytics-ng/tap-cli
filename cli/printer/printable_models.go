package printer

import (
	"strconv"
	"strings"
	"time"

	apiServiceModels "github.com/trustedanalytics/tap-api-service/models"
	userManagement "github.com/trustedanalytics/tap-api-service/user-management-connector"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-cli/api"
)

const timeFormatter = "Jan 02 15:04"
const LastMessageMark = "..."

type Printable interface {
	Headers() []string
	StandarizedData() []string
}

type PrintableOffering struct {
	apiServiceModels.Offering
}

func (po PrintableOffering) Headers() []string {
	return []string{"name", "plan", "description", "state"}
}
func (po PrintableOffering) StandarizedData() []string {
	planNames := []string{}
	for _, planName := range po.OfferingPlans {
		planNames = append(planNames, planName.Name)
	}
	return []string{po.Name, strings.Join(planNames, ", "), po.Description, po.State}
}

type PrintableService struct {
	apiServiceModels.ServiceInstance
}

func (s PrintableService) Headers() []string {
	return []string{"name", "service", "plan", "state", "created by", "create", "updated by", "update", "message"}
}
func (s PrintableService) StandarizedData() []string {
	return []string{s.Name, s.ServiceName, s.ServicePlanName, s.State.String(),
		s.AuditTrail.CreatedBy, formatTime(s.AuditTrail.CreatedOn),
		s.AuditTrail.LastUpdateBy, formatTime(s.AuditTrail.LastUpdatedOn),
		getLastMessageMark(s.Metadata)}
}

type PrintableApplication struct {
	apiServiceModels.ApplicationInstance
}

func (app PrintableApplication) Headers() []string {
	return []string{"name", "image state", "state", "replication", "memory", "disk", "urls", "created by", "create", "updated by", "update", "message"}
}
func (app PrintableApplication) StandarizedData() []string {
	return []string{app.Name, string(app.ImageState), app.State.String(),
		strconv.Itoa(app.Replication), app.Memory, app.DiskQuota, strings.Join(app.Urls, ","),
		app.AuditTrail.CreatedBy, formatTime(app.AuditTrail.CreatedOn),
		app.AuditTrail.LastUpdateBy, formatTime(app.AuditTrail.LastUpdatedOn),
		getLastMessageMark(app.Metadata)}
}

type PrintableRecentlyPushedApplication struct {
	catalogModels.Application
}

func (app PrintableRecentlyPushedApplication) Headers() []string {
	return []string{"name", "image id", "description", "replication", "created by", "create", "updated by", "update"}
}
func (app PrintableRecentlyPushedApplication) StandarizedData() []string {
	return []string{app.Name, app.ImageId, app.Description,
		strconv.Itoa(app.Replication), app.AuditTrail.CreatedBy, formatTime(app.AuditTrail.CreatedOn),
		app.AuditTrail.LastUpdateBy, formatTime(app.AuditTrail.LastUpdatedOn)}
}

type PrintableCredentials struct {
	api.Credentials
}

func (pc PrintableCredentials) Headers() []string {
	return []string{"api", "username"}
}
func (pc PrintableCredentials) StandarizedData() []string {
	return []string{pc.Address, pc.Username}
}

type PrintableUser struct {
	userManagement.UaaUser
}

func (pu PrintableUser) Headers() []string {
	return []string{"username"}
}
func (pu PrintableUser) StandarizedData() []string {
	return []string{pu.Username}
}

type PrintableInvitation struct {
	Email string
}

func (pi PrintableInvitation) Headers() []string {
	return []string{"e-mail"}
}
func (pi PrintableInvitation) StandarizedData() []string {
	return []string{pi.Email}
}

type PrintableResource struct {
	apiServiceModels.InstanceBindingsResource
}

func (pb PrintableResource) Headers() []string {
	return []string{"binding name", "binding id"}
}
func (pb PrintableResource) StandarizedData() []string {
	return []string{pb.ServiceInstanceName, pb.ServiceInstanceGuid}
}

func getLastMessageMark(metadata []catalogModels.Metadata) string {
	if catalogModels.GetValueFromMetadata(metadata, catalogModels.LAST_STATE_CHANGE_REASON) != "" {
		return LastMessageMark
	}
	return ""
}

func formatTime(t int64) string {
	return time.Unix(t, 0).Format(timeFormatter)
}
