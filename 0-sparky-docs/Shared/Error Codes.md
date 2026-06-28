All error codes are to be returned via `X-Sparky-ErrorCode`, as well as in an HTTP body, in which the json struct contains two root elements:
- `ErrorCode` - the code in question
- `ErrorDetail` - why that code may've been hit.

#Retriable errors will always come with a `Retry-After` header. 
# Codes
# Model loading errors
### ModelNotRegistered
The requested model is not present in registrations.
### ModelTooLarge
The requested model requires more nodes than are physically available in the cluster.
### CannotEvict #Retriable 
The model is large enough to live on the cluster, but the cluster is too busy to automatically evict.
### ModelLoading #Retriable
This model is currently loading.