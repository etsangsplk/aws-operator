package cloudformation

import (
	"context"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/framework"

	"github.com/giantswarm/aws-operator/service/key"
)

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	// no-op if we are not using cloudformation
	if !key.UseCloudFormation(customObject) {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "not processing cloudformation stack")
		return nil
	}

	updateStackInput, err := toUpdateStackInput(updateChange)
	if err != nil {
		return microerror.Mask(err)
	}

	stackName := updateStackInput.StackName
	if *stackName != "" {
		_, err := r.awsClient.UpdateStack(&updateStackInput)
		if err != nil {
			return microerror.Maskf(err, "updating AWS cloudformation stack")
		}

		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "updating AWS cloudformation stack: updated")
	} else {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "updating AWS cloudformation stack: no need to update")
	}

	return nil
}

func (r *Resource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	create, err := r.newCreateChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	update, err := r.newUpdateChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := framework.NewPatch()
	patch.SetCreateChange(create)
	patch.SetUpdateChange(update)

	return patch, nil
}

func (r *Resource) newUpdateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return awscloudformation.CreateStackInput{}, microerror.Mask(err)
	}

	desiredStackState, err := toStackState(desiredState)
	if err != nil {
		return awscloudformation.CreateStackInput{}, microerror.Mask(err)
	}

	currentStackState, err := toStackState(currentState)
	if err != nil {
		return awscloudformation.CreateStackInput{}, microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "finding out if the main stack should be updated")

	updateState := awscloudformation.UpdateStackInput{
		StackName: aws.String(""),
	}

	if !reflect.DeepEqual(desiredStackState, currentStackState) {
		var mainTemplate string
		/*
			      commented out until we assing proper values to the template
						mainTemplate, err := getMainTemplateBody(customObject)
						if err != nil {
							return nil, microerror.Mask(err)
						}
		*/
		updateState.StackName = aws.String(desiredStackState.Name)
		updateState.TemplateBody = aws.String(mainTemplate)
		updateState.Parameters = []*awscloudformation.Parameter{
			{
				ParameterKey:   aws.String(workersParameterKey),
				ParameterValue: aws.String(desiredStackState.Workers),
			},
			{
				ParameterKey:   aws.String(imageIDParameterKey),
				ParameterValue: aws.String(desiredStackState.ImageID),
			},
			{
				ParameterKey:   aws.String(clusterVersionParameterKey),
				ParameterValue: aws.String(desiredStackState.ClusterVersion),
			},
		}
	}

	return updateState, nil
}