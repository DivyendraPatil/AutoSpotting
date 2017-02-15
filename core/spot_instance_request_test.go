package autospotting

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type mock struct {
	ec2iface.EC2API
	er      error
	dsiro   *ec2.DescribeSpotInstanceRequestsOutput
	dsiroer error
	dio     *ec2.DescribeInstancesOutput
	dioer   error
}

func (m mock) CreateTags(in *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return nil, m.er
}

func (m mock) WaitUntilSpotInstanceRequestFulfilled(in *ec2.DescribeSpotInstanceRequestsInput) error {
	return m.er
}

func (m mock) DescribeSpotInstanceRequests(in *ec2.DescribeSpotInstanceRequestsInput) (*ec2.DescribeSpotInstanceRequestsOutput, error) {
	return m.dsiro, m.dsiroer
}

func (m mock) DescribeInstances(in *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return m.dio, m.dioer
}

func Test_waitForAndTagSpotInstance(t *testing.T) {
	tests := []struct {
		name string
		req  spotInstanceRequest
		er   error
	}{
		{
			name: "with WaitUntilSpotInstanceRequestFulfilled error",
			req: spotInstanceRequest{
				SpotInstanceRequest: &ec2.SpotInstanceRequest{
					SpotInstanceRequestId: aws.String(""),
				},
				region: &region{
					services: connections{
						ec2: mock{
							er: errors.New(""),
						},
					},
				},
				asg: &autoScalingGroup{
					name: ""},
			},
			er: errors.New(""),
		},
		{
			name: "without WaitUntilSpotInstanceRequestFulfilled error",
			req: spotInstanceRequest{
				SpotInstanceRequest: &ec2.SpotInstanceRequest{
					SpotInstanceRequestId: aws.String(""),
				},
				region: &region{
					services: connections{
						ec2: mock{
							dsiro: &ec2.DescribeSpotInstanceRequestsOutput{
								SpotInstanceRequests: []*ec2.SpotInstanceRequest{
									{InstanceId: aws.String("")},
								},
							},
							dio: &ec2.DescribeInstancesOutput{
								Reservations: []*ec2.Reservation{{}},
							},
						},
					},
				},
				asg: &autoScalingGroup{
					Group: &autoscaling.Group{
						Tags: []*autoscaling.TagDescription{},
					},
					name: "",
				},
			},
			er: errors.New(""),
		},
		{
			name: "with DescribeSpotInstanceRequestsOutput error",
			req: spotInstanceRequest{
				SpotInstanceRequest: &ec2.SpotInstanceRequest{
					SpotInstanceRequestId: aws.String(""),
				},
				region: &region{
					services: connections{
						ec2: mock{
							dsiro: &ec2.DescribeSpotInstanceRequestsOutput{
								SpotInstanceRequests: []*ec2.SpotInstanceRequest{
									{InstanceId: aws.String("")},
								},
							},
							dsiroer: errors.New(""),
							dio: &ec2.DescribeInstancesOutput{
								Reservations: []*ec2.Reservation{{}},
							},
						},
					},
				},
				asg: &autoScalingGroup{
					Group: &autoscaling.Group{
						Tags: []*autoscaling.TagDescription{},
					},
					name: "",
				},
			},
			er: errors.New(""),
		},
	}

	for _, tc := range tests {
		tc.req.waitForAndTagSpotInstance()
	}
}

func Test_tag(t *testing.T) {
	tests := []struct {
		name string
		tag  string
		req  spotInstanceRequest
		er   error
	}{
		{
			name: "with error",
			tag:  "tag",
			req: spotInstanceRequest{
				SpotInstanceRequest: &ec2.SpotInstanceRequest{
					SpotInstanceRequestId: aws.String(""),
				},
				region: &region{
					services: connections{
						ec2: mock{
							er: errors.New(""),
						},
					},
				},
			},
			er: errors.New(""),
		},
		{
			name: "without error",
			tag:  "tag",
			req: spotInstanceRequest{
				SpotInstanceRequest: &ec2.SpotInstanceRequest{
					SpotInstanceRequestId: aws.String(""),
				},
				region: &region{
					services: connections{
						ec2: mock{
							er: nil,
						},
					},
				},
			},
			er: nil,
		},
	}

	for _, tc := range tests {
		er := tc.req.tag(tc.tag)
		if er != nil && er.Error() != tc.er.Error() {
			t.Errorf("error actual: %s, expected: %s", er.Error(), tc.er.Error())
		}
	}
}
