#!/bin/bash

# Variables you may need to edit
NUM_OF_INSTANCES=3
INIT_SCRIPT="./init-agent-box.sh"
TARGET_SECRET=$THINGZ_SECRET
TARGET_HOST=$THINGZ_HOST
AGENT_STRGY="cpu:3,cpus:60,mem:3,swap:30,load:5"
AWS_KEY="my-key-rsa" # name of the key on aws, should be in ~/.ssh/ too

# Other AWS variables
AWS_PROFILE="admin" # aws cli profile in ~/.aws/config
AWS_SEC_GROUP="thingz-agent-sec-group" # name of the EC2 sec group to use
AWS_REGION="us-west-2" # region to deploy too
AWS_AMI="ami-838dd9b3" # base Ubuntu image to use (14.04)
AWS_INSTANCE="t2.micro" # instance type
AWS_TAG="thingz-agent" # tag name for ease of deletes later, will append index to


echo "Deploying thingz-agent instances: $NUM_OF_INSTANCES"
export iids=$(aws ec2 run-instances \
                        --count $NUM_OF_INSTANCES \
                        --image-id $AWS_AMI \
                        --key-name $AWS_KEY \
                        --security-groups $AWS_SEC_GROUP \
                        --instance-type $AWS_INSTANCE \
                        --region $AWS_REGION \
                        --profile $AWS_PROFILE \
                        --output text \
                        --query 'Instances[*].InstanceId')

echo "Provisioned instances: $iids"

iindex=0

for iid in $iids
do
    echo "Waiting for instance $iindex ($iid) to start running:"
    while state=$(aws ec2 describe-instances \
                        --instance-ids $iid \
                        --region $AWS_REGION \
                        --profile $AWS_PROFILE \
                        --output text \
                        --query 'Reservations[*].Instances[*].State.Name'); test "$state" != "running"; do
        sleep 1; echo -n "."
    done;
    echo -e "\nInstance $iindex is $state"

    echo "Tagging instance $iid: name=thingz-agent-$iindex"
    itag=$(aws ec2 create-tags \
                        --resources $iid \
                        --region $AWS_REGION \
                        --profile $AWS_PROFILE \
                        --tags Key=name,Value=$AWS_TAG-$iindex)

    echo "Getting instance $iid public IP"
    iip=$(aws ec2 describe-instances \
                        --instance-ids $iid \
                        --region $AWS_REGION \
                        --profile $AWS_PROFILE \
                        --output text \
                        --query 'Reservations[*].Instances[*].PublicIpAddress')
    echo "Instance public IP: $iip"

    echo "Configuring instance $iid"
    ssh -i ~/.ssh/$AWS_KEY \
        -o StrictHostKeyChecking=no \
        -o UserKnownHostsFile=/dev/null \
        -o ConnectTimeout=60 \
        -o ConnectionAttempts=3 \
        ubuntu@$iip 'bash -s' < $INIT_SCRIPT

    echo "Coping thingz-agent assets to instance: $iindex"
    rsync -i ~/.ssh/$AWS_KEY -r -a  \
          -e "ssh -ax -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null" \
          --exclude '.DS_Store' \
          --exclude '.git' \
          --exclude 'scripts' \
          --exclude 'thingz-agent' \
          $GOPATH/src/github.com/mchmarny/thingz-agent/ \
          ubuntu@$iip:/home/ubuntu/go/src/github.com/mchmarny/thingz-agent

    echo "Coping thingz-agent dependencies to instance: $iindex"
    rsync -i ~/.ssh/$AWS_KEY -r -a  \
          -e "ssh -ax -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null" \
          --exclude '.DS_Store' \
          --exclude '.git' \
          $GOPATH/src/github.com/mchmarny/thingz-commons/ \
          ubuntu@$iip:/home/ubuntu/go/src/github.com/mchmarny/thingz-commons

    echo "Ready for SSH:"
    echo "ssh -i ~/.ssh/$AWS_KEY ubuntu@$iip"

    echo "Starting agent on $iip ($iid)"
    ssh -i ~/.ssh/$AWS_KEY \
        -o StrictHostKeyChecking=no \
        -o UserKnownHostsFile=/dev/null \
        -o ConnectTimeout=60 \
        -o ConnectionAttempts=3 \
        ubuntu@$iip "source /etc/profile; cd \$GOPATH/src/github.com/mchmarny/thingz-agent; go build; echo Launching agent: \$HOSTNAME; nohup ./thingz-agent --strategy=$AGENT_STRGY --publisher=influxdb --publisher-args=udp://agent:${TARGET_SECRET}@${TARGET_HOST}:4444/thingz > ./thingz.log 2>&1 &"

    ((iindex++))

done