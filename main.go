package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		sgArgs := &ec2.SecurityGroupArgs{
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol: pulumi.String("tcp"),
					FromPort: pulumi.Int(8080),
					ToPort: pulumi.Int(8080),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
				ec2.SecurityGroupIngressArgs{
					Protocol: pulumi.String("tcp"),
					FromPort: pulumi.Int(22),
					ToPort: pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Protocol: pulumi.String("-1"),
					FromPort: pulumi.Int(0),
					ToPort: pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},

				},
			},
		}
		sg,err := ec2.NewSecurityGroup(ctx,"jenkins-sg",sgArgs)
		if err != nil{
			return err
		}
		kp,err := ec2.NewKeyPair(ctx,"local-ssh",&ec2.KeyPairArgs{
			PublicKey: pulumi.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDcg00wCdknCnflaBjqt2/g87DFElzlC+EAK8PTVWMgtqDgwEcFBaUyghE8M/gTezlz4S+pL7HkQkkxz2IiWAdoe8wj0xWTUuy9nei5EOV33ZAKj/8v0OvVobtCMVBB1nrK0lLr5vmfSNT3NFbUcsWx0IkD8CqcNLNIbCZfvW1T/D3WXWMV6IpD8Gwl8HrEylaxx7Vrj3SlstPlNgHfqSksf7kE3HZvI8Br09H+b3FPF5Dp+IzeeRl1EGVWNlTAIObC4OfAMtYLEyxUp+DdhiTTn1n6YrHb250X1vjMBIzsBzuVs1e9C5lf3a/pTVR4CLpU40ulJk5KmUL9t07Qxxr3B1Mh6/fFK6cVJato4kB4ScU8E/F+EQyrTnYdOP9onwc4LmD4Tp+H7V5smwepaE1bUpdT03NQ9R+LdiK1OHvL2n9E/7ETBHWaGSUawI7JfB7GOE3r+wboaXruD+7AEo/sTIzO3M3x0thPlCRawfxMYcmkgPXsFhCmWAS7nrVtkyU= khomsankhuarkhachon@Khomsans-MacBook-Pro.local"),
		})
		if err != nil {
			return nil
		}
		jenkinsServer,err := ec2.NewInstance(ctx,"jenkins-server",&ec2.InstanceArgs{
			InstanceType:pulumi.String("t2.micro"),
			VpcSecurityGroupIds: pulumi.StringArray{sg.ID()},
			Ami: pulumi.String("ami-0896ef7bec0d0e792"),
			KeyName: kp.KeyName,
		})
		if err != nil {
			return nil
		}
		fmt.Println(jenkinsServer.PublicIp)
		fmt.Println(jenkinsServer.PublicDns)

		ctx.Export("publicIp",jenkinsServer.PublicIp)
		ctx.Export("publicHostname",jenkinsServer.PublicDns)

		return nil
	})
}
