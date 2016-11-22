Kublambda - lambda for kubernetes 

Features: 
- new third party resource called "lambda" 
    - name of lambda function 
    - source code of lambda function 
    - events that should trigger lambda function 
        - http
        - resource watch 
    - names of state tables for lambda 

Third party resource 'lambda' 
-----------------------------
When kublambda is installed, it creates a new resource called "lambda". It also starts a kublambda controller that manages the lambda tpr lifecycle. To create a lambda tpr, the user must specify
- the name of the lambda
- the source code of the lambda
- the events that can trigger the lambda 

The controller takes care of processing the source code to turn it into an executable asset. This executable is stored in an object store.

Calling a Lambda 
----------------
A request arrives at hte runner service with some URI. The service forwards the request to an available runner. The runner uses the PATH portion of the request as a key to lookup into the source code store. The resulting binary is executed for the request, being passed the original http request. 

Horizontal autoscaling
--------------------------
The system runs a Deployment of lambda runners. An autoscaler looks over the pool of runner pods, and auto-scales the replicaset of runners in response to load. 

HTTP trigger lambdas
--------------------
If HTTP is enabled as a trigger, the controller will create an ingress resource for the service backing the function. This will cause the service to be exposed outside of the kubernetes cluster. The URL of the lambda is stored as a read-only property in the lambda resource. 

Resource event triggering of lambdas
----------------------------------
If a resource event is specified for a lambda, the controller will launch a pod that implements a watch-loop on the specified resource stream. For each event in the stream, it will evaluate the condition specified in the lambda tpr and, if it matches, call the lambda's service. 


Lambda image management
-----------------------
The Lambda controller automatically builds a container image for the user code and pushes it to the container registry. The lambda controller can automatically delete these images when lambdas are deleted. The image uri is stored as a read-only property in the lambda tpr. 
