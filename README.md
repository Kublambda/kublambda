Kublambda - lambda for kubernetes 

Features: 
- new third party resource called "lambda" 
-- name of lambda function 
-- source code of lambda function 
-- events that should trigger lambda function 
---- http
---- resource watch 
-- names of state tables for lambda 

Third party resource 'lambda' 
-----------------------------
When kublambda is installed, it creates a new resource called "lambda". It also starts a kublambda controller that manages the lambda tpr lifecycle. To create a lambda tpr, the user must specify
- the name of the lambda
- the source code of the lambda
- the events that can trigger the lambda 

The controller takes care of processing the source code to turn it into an executable asset. This executable is stored in a registry (another tpr?).

Horizontal autoscaling
--------------------------
When a lambda is created, the controller will create a service and a deployment and autoscaler for the lambda, using a docker-containerized version of the user's source code. At this point, autoscaler requires the minimum node count to be 1 (or greater), so there will always be at least one instance of the function, even if it's idle. But the k8s docs state that a future goal is to allow the minimal number of pods to be 0. This setup allows the system to allocate more pods to the lambda if load increases, and remove that allocation as load falls.   

HTTP trigger lambdas
--------------------
If HTTP is enabled as a trigger, the controller will create an ingress resource, and or service-type (nodeport or load balancer) for the service backing the function. This will cause the service to be exposed outside of the kubernetes cluster. The URL of the lambda is stored as a read-only property in the lambda resource. 

Resource event triggering of lambdas
----------------------------------
If a resource event is specified for a lambda, the controller will launch a pod that implements a watch-loop on the specified resource stream. For each event in the stream, it will evaluate the condition specified in the lambda tpr and, if it matches, call the lambda's service. 


Lambda image management
-----------------------
The Lambda controller automatically builds a container image for the user code and pushes it to the container registry. The lambda controller can automatically delete these images when lambdas are deleted. The image uri is stored as a read-only property in the lambda tpr. 



scenario: create a new lambda tpr 
-----------------------------
inputs:
- name
- code (file or zipfile) 
- triggers 

procedure: 
1. launch a pod to build the source code 
2. executes buildpack for code
3. stores binary in storage somewhere? 
4. creates service, deployment & autoscaler for function
5. if http trigger specified, create ingres for function
6. if event trigger specified, create event-controller for function 

scneario: change code in running lambda
---------------------------------------
procedure: 
1. launch a pod to build the source code
2. executes buildpack for code 
3. stores binary in storage somewhere? 
4. updates deployment to new image name 
5. rollout change to deployment 

design thoughts: 
- instead of running a deployment per function, we can have one big deployment
- each pod in the deployment runs a runtime that can execute any lambda
- request arrives at any pod in cluster
- pod checks to see if it has code for function in cache 
- if not, fetch code from object storage store it in cache
- run user code in a sandbox 
- return results to caller 

