package ovirt

import (
	"time"

	"github.com/go-logr/logr"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
	vmv1alpha1 "github.com/tmax-cloud/hypercloud-ovirt-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	Timeout = 10 * time.Second
)

// OvirtActuator contains connection data
type OvirtActuator struct {
	conn *ovirtsdk4.Connection
	url  string
	name string
	pass string
}

// NewActuator creates new OvirtActuator
func NewActuator() *OvirtActuator {
	return &OvirtActuator{
		conn: nil,
		url:  "https://node1.test.dom/ovirt-engine/api",
		name: "admin@internal",
		pass: "1",
	}
}

func (actuator *OvirtActuator) getConnection() (*ovirtsdk4.Connection, error) {
	var err error
	if actuator.conn == nil || actuator.conn.Test() != nil {
		actuator.conn, err = actuator.createConnection()
	}

	return actuator.conn, err
}

func (actuator *OvirtActuator) createConnection() (*ovirtsdk4.Connection, error) {
	conn, err := ovirtsdk4.NewConnectionBuilder().
		URL(actuator.url).
		Username(actuator.name).
		Password(actuator.pass).
		Insecure(true).
		Compress(true).
		Timeout(Timeout).
		Build()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// AddVM gets the virtual machine from Ovirt cluster
func (actuator *OvirtActuator) GetVM(log logr.Logger, m *vmv1alpha1.VirtualMachine) error {
	conn, err := actuator.getConnection()
	if err != nil {
		log.Error(err, "Make connection failed")
		return err
	}
	defer conn.Close()

	vmsService := conn.SystemService().VmsService()
	vmsResponse, err := vmsService.List().Search("name=" + m.Name).Send()
	if err != nil {
		log.Error(err, "Failed to get vm list")
		return err
	}
	vms, _ := vmsResponse.Vms()
	if vm := vms.Slice(); vm != nil {
		return nil
	}

	return errors.NewNotFound(schema.GroupResource{}, m.Name)
}

// AddVM adds the virtual machine to Ovirt cluster
func (actuator *OvirtActuator) AddVM(log logr.Logger, m *vmv1alpha1.VirtualMachine) error {
	conn, err := actuator.getConnection()
	if err != nil {
		log.Error(err, "Make connection failed")
		return err
	}
	defer conn.Close()

	vmsService := conn.SystemService().VmsService()
	cluster, err := ovirtsdk4.NewClusterBuilder().Name("Default").Build()
	if err != nil {
		log.Error(err, "Failed to build cluster")
		return err
	}
	if m.Spec.Template == "" {
		m.Spec.Template = "Blank"
	}
	template, err := ovirtsdk4.NewTemplateBuilder().Name(m.Spec.Template).Build()
	if err != nil {
		log.Error(err, "Failed to build template")
		return err
	}
	vm, err := ovirtsdk4.NewVmBuilder().Name(m.Name).Cluster(cluster).Template(template).Build()
	if err != nil {
		log.Error(err, "Failed to build vm")
		return err
	}
	resp, err := vmsService.Add().Vm(vm).Send()
	if err != nil {
		log.Error(err, "Failed to add vm")
		return err
	}

	vm, _ = resp.Vm()
	name, _ := vm.Name()
	log.Info("Add vm successfully", "vm.Name", name)

	return nil
}

// FinalizeVm removes the virtual machine from Ovirt cluster
func (actuator *OvirtActuator) FinalizeVm(log logr.Logger, m *vmv1alpha1.VirtualMachine) error {
	conn, err := actuator.getConnection()
	if err != nil {
		log.Error(err, "Make connection failed")
		return err
	}
	defer conn.Close()

	vmsService := conn.SystemService().VmsService()
	vmsResponse, err := vmsService.List().Search("name=" + m.Name).Send()
	if err != nil {
		log.Error(err, "Failed to search vms")
		return err
	}
	vms, _ := vmsResponse.Vms()
	id, _ := vms.Slice()[0].Id()
	vmService := vmsService.VmService(id)
	_, err = vmService.Remove().Send()
	if err != nil {
		log.Error(err, "Failed to remove vm")
		return err
	}

	log.Info("Remove vm successfully", "vm.Name", m.Name)
	return nil
}