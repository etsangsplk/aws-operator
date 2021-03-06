package cloudformation

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/giantswarm/aws-operator/service/awsconfig/v3/key"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/framework"
)

func (r *Resource) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	cluster, err := key.ToCustomObject(obj)
	if err != nil {
		microerror.Mask(err)
	}

	deleteStackInput := awscloudformation.DeleteStackInput{
		StackName: aws.String(key.MainGuestStackName(cluster)),
	}
	_, err = r.Clients.CloudFormation.DeleteStack(&deleteStackInput)
	if err != nil {
		return microerror.Maskf(err, "deleting AWS Guest CloudFormation Stack")
	}
	r.logger.LogCtx(ctx, "debug", "deleting AWS Guest CloudFormation stack: deleted")

	deleteStackInput = awscloudformation.DeleteStackInput{
		StackName: aws.String(key.MainHostPreStackName(cluster)),
	}
	_, err = r.HostClients.CloudFormation.DeleteStack(&deleteStackInput)
	if err != nil {
		return microerror.Maskf(err, "deleting AWS Host Pre-Guest CloudFormation Stack")
	}
	r.logger.LogCtx(ctx, "debug", "deleting AWS Host Pre-Guest CloudFormation stack: deleted")

	deleteStackInput = awscloudformation.DeleteStackInput{
		StackName: aws.String(key.MainHostPostStackName(cluster)),
	}
	_, err = r.HostClients.CloudFormation.DeleteStack(&deleteStackInput)
	if err != nil {
		return microerror.Maskf(err, "deleting AWS Host Post-Guest CloudFormation Stack")
	}
	r.logger.LogCtx(ctx, "debug", "deleting AWS Host Post-Guest CloudFormation stack: deleted")

	return nil
}

func (r *Resource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	delete, err := r.newDeleteChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := framework.NewPatch()
	patch.SetDeleteChange(delete)

	return patch, nil
}

func (r *Resource) newDeleteChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentStackState, err := toStackState(currentState)
	if err != nil {
		return awscloudformation.DeleteStackInput{}, microerror.Mask(err)
	}

	desiredStackState, err := toStackState(desiredState)
	if err != nil {
		return awscloudformation.DeleteStackInput{}, microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "debug", "finding out if the main stack should be deleted")

	deleteState := awscloudformation.DeleteStackInput{
		StackName: aws.String(""),
	}

	if currentStackState.Name != "" && desiredStackState.Name != currentStackState.Name {
		deleteState.StackName = aws.String(currentStackState.Name)
	}

	return deleteState, nil
}
