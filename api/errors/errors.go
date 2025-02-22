package errors

import (
	"fmt"
	"net/http"
)

type (
	// ErrorCode type of error status.
	ErrorCode int

	errorCodeMap map[ErrorCode]Error

	// Error structure represents API error.
	Error struct {
		ErrCode        ErrorCode
		Code           string
		Description    string
		HTTPStatusCode int
	}
)

const maxEConfigJSONSize = 262272

// Error codes, non exhaustive list - http://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html
const (
	_ ErrorCode = iota
	ErrAccessDenied
	ErrBadDigest
	ErrEntityTooSmall
	ErrEntityTooLarge
	ErrPolicyTooLarge
	ErrIllegalVersioningConfigurationException
	ErrIncompleteBody
	ErrInternalError
	ErrInvalidAccessKeyID
	ErrInvalidBucketName
	ErrInvalidDigest
	ErrInvalidRange
	ErrInvalidCopyPartRange
	ErrInvalidCopyPartRangeSource
	ErrInvalidMaxKeys
	ErrInvalidEncodingMethod
	ErrInvalidMaxUploads
	ErrInvalidMaxParts
	ErrInvalidPartNumberMarker
	ErrInvalidRequestBody
	ErrInvalidCopySource
	ErrInvalidMetadataDirective
	ErrInvalidCopyDest
	ErrInvalidPolicyDocument
	ErrInvalidObjectState
	ErrMalformedXML
	ErrMissingContentLength
	ErrMissingContentMD5
	ErrMissingRequestBodyError
	ErrMissingSecurityHeader
	ErrNoSuchBucket
	ErrNoSuchBucketPolicy
	ErrNoSuchBucketLifecycle
	ErrNoSuchLifecycleConfiguration
	ErrNoSuchBucketSSEConfig
	ErrNoSuchCORSConfiguration
	ErrNoSuchWebsiteConfiguration
	ErrReplicationConfigurationNotFoundError
	ErrNoSuchKey
	ErrNoSuchUpload
	ErrNoSuchVersion
	ErrInvalidVersion
	ErrInvalidArgument
	ErrInvalidTagKey
	ErrInvalidTagValue
	ErrInvalidTagsSizeExceed
	ErrNotImplemented
	ErrPreconditionFailed
	ErrNotModified
	ErrRequestTimeTooSkewed
	ErrSignatureDoesNotMatch
	ErrMethodNotAllowed
	ErrInvalidPart
	ErrInvalidPartOrder
	ErrAuthorizationHeaderMalformed
	ErrMalformedPOSTRequest
	ErrPOSTFileRequired
	ErrSignatureVersionNotSupported
	ErrBucketNotEmpty
	ErrAllAccessDisabled
	ErrMalformedPolicy
	ErrMissingFields
	ErrMissingCredTag
	ErrCredMalformed
	ErrInvalidRegion
	ErrInvalidServiceS3
	ErrInvalidServiceSTS
	ErrInvalidRequestVersion
	ErrMissingSignTag
	ErrMissingSignHeadersTag
	ErrMalformedDate
	ErrMalformedPresignedDate
	ErrMalformedCredentialDate
	ErrMalformedCredentialRegion
	ErrMalformedExpires
	ErrNegativeExpires
	ErrAuthHeaderEmpty
	ErrExpiredPresignRequest
	ErrRequestNotReadyYet
	ErrUnsignedHeaders
	ErrMissingDateHeader
	ErrInvalidQuerySignatureAlgo
	ErrInvalidQueryParams
	ErrBucketAlreadyOwnedByYou
	ErrInvalidDuration
	ErrBucketAlreadyExists
	ErrMetadataTooLarge
	ErrUnsupportedMetadata
	ErrMaximumExpires
	ErrSlowDown
	ErrInvalidPrefixMarker
	ErrBadRequest
	ErrKeyTooLongError
	ErrInvalidBucketObjectLockConfiguration
	ErrObjectLockConfigurationNotFound
	ErrObjectLockConfigurationNotAllowed
	ErrNoSuchObjectLockConfiguration
	ErrObjectLocked
	ErrInvalidRetentionDate
	ErrPastObjectLockRetainDate
	ErrUnknownWORMModeDirective
	ErrBucketTaggingNotFound
	ErrObjectLockInvalidHeaders
	ErrInvalidTagDirective
	// Add new error codes here.
	ErrNotSupported

	// SSE-S3 related API errors.
	ErrInvalidEncryptionMethod

	// Server-Side-Encryption (with Customer provided key) related API errors.
	ErrInsecureSSECustomerRequest
	ErrSSEMultipartEncrypted
	ErrSSEEncryptedObject
	ErrInvalidEncryptionParameters
	ErrInvalidSSECustomerAlgorithm
	ErrInvalidSSECustomerKey
	ErrMissingSSECustomerKey
	ErrMissingSSECustomerKeyMD5
	ErrSSECustomerKeyMD5Mismatch
	ErrInvalidSSECustomerParameters
	ErrIncompatibleEncryptionMethod
	ErrKMSNotConfigured
	ErrKMSAuthFailure

	ErrNoAccessKey
	ErrInvalidToken

	// Bucket notification related errors.
	ErrEventNotification
	ErrARNNotification
	ErrRegionNotification
	ErrOverlappingFilterNotification
	ErrFilterNameInvalid
	ErrFilterNamePrefix
	ErrFilterNameSuffix
	ErrFilterValueInvalid
	ErrOverlappingConfigs
	ErrUnsupportedNotification

	// S3 extended errors.
	ErrContentSHA256Mismatch

	// Add new extended error codes here.

	// MinIO extended errors.
	//   ErrReadQuorum
	//   ErrWriteQuorum
	ErrParentIsObject
	ErrStorageFull
	ErrRequestBodyParse
	ErrObjectExistsAsDirectory
	ErrInvalidObjectName
	ErrInvalidObjectNamePrefixSlash
	ErrInvalidResourceName
	ErrServerNotInitialized
	ErrOperationTimedOut
	ErrOperationMaxedOut
	ErrInvalidRequest
	// MinIO storage class error codes.
	ErrInvalidStorageClass
	ErrBackendDown
	// Add new extended error codes here.
	// Please open a https://github.com/minio/minio/issues before adding
	// new error codes here.

	ErrMalformedJSON
	ErrAdminNoSuchUser
	ErrAdminNoSuchGroup
	ErrAdminGroupNotEmpty
	ErrAdminNoSuchPolicy
	ErrAdminInvalidArgument
	ErrAdminInvalidAccessKey
	ErrAdminInvalidSecretKey
	ErrAdminConfigNoQuorum
	ErrAdminConfigTooLarge
	ErrAdminConfigBadJSON
	ErrAdminConfigDuplicateKeys
	ErrAdminCredentialsMismatch
	ErrInsecureClientRequest
	ErrObjectTampered
	// Bucket Quota error codes.
	ErrAdminBucketQuotaExceeded
	ErrAdminNoSuchQuotaConfiguration
	ErrAdminBucketQuotaDisabled

	ErrHealNotImplemented
	ErrHealNoSuchProcess
	ErrHealInvalidClientToken
	ErrHealMissingBucket
	ErrHealAlreadyRunning
	ErrHealOverlappingPaths
	ErrIncorrectContinuationToken

	// S3 Select Errors.
	ErrEmptyRequestBody
	ErrUnsupportedFunction
	ErrInvalidExpressionType
	ErrBusy
	ErrUnauthorizedAccess
	ErrExpressionTooLong
	ErrIllegalSQLFunctionArgument
	ErrInvalidKeyPath
	ErrInvalidCompressionFormat
	ErrInvalidFileHeaderInfo
	ErrInvalidJSONType
	ErrInvalidQuoteFields
	ErrInvalidRequestParameter
	ErrInvalidDataType
	ErrInvalidTextEncoding
	ErrInvalidDataSource
	ErrInvalidTableAlias
	ErrMissingRequiredParameter
	ErrObjectSerializationConflict
	ErrUnsupportedSQLOperation
	ErrUnsupportedSQLStructure
	ErrUnsupportedSyntax
	ErrUnsupportedRangeHeader
	ErrLexerInvalidChar
	ErrLexerInvalidOperator
	ErrLexerInvalidLiteral
	ErrLexerInvalidIONLiteral
	ErrParseExpectedDatePart
	ErrParseExpectedKeyword
	ErrParseExpectedTokenType
	ErrParseExpected2TokenTypes
	ErrParseExpectedNumber
	ErrParseExpectedRightParenBuiltinFunctionCall
	ErrParseExpectedTypeName
	ErrParseExpectedWhenClause
	ErrParseUnsupportedToken
	ErrParseUnsupportedLiteralsGroupBy
	ErrParseExpectedMember
	ErrParseUnsupportedSelect
	ErrParseUnsupportedCase
	ErrParseUnsupportedCaseClause
	ErrParseUnsupportedAlias
	ErrParseUnsupportedSyntax
	ErrParseUnknownOperator
	ErrParseMissingIdentAfterAt
	ErrParseUnexpectedOperator
	ErrParseUnexpectedTerm
	ErrParseUnexpectedToken
	ErrParseUnexpectedKeyword
	ErrParseExpectedExpression
	ErrParseExpectedLeftParenAfterCast
	ErrParseExpectedLeftParenValueConstructor
	ErrParseExpectedLeftParenBuiltinFunctionCall
	ErrParseExpectedArgumentDelimiter
	ErrParseCastArity
	ErrParseInvalidTypeParam
	ErrParseEmptySelect
	ErrParseSelectMissingFrom
	ErrParseExpectedIdentForGroupName
	ErrParseExpectedIdentForAlias
	ErrParseUnsupportedCallWithStar
	ErrParseNonUnaryAgregateFunctionCall
	ErrParseMalformedJoin
	ErrParseExpectedIdentForAt
	ErrParseAsteriskIsNotAloneInSelectList
	ErrParseCannotMixSqbAndWildcardInSelectList
	ErrParseInvalidContextForWildcardInSelectList
	ErrIncorrectSQLFunctionArgumentType
	ErrValueParseFailure
	ErrEvaluatorInvalidArguments
	ErrIntegerOverflow
	ErrLikeInvalidInputs
	ErrCastFailed
	ErrInvalidCast
	ErrEvaluatorInvalidTimestampFormatPattern
	ErrEvaluatorInvalidTimestampFormatPatternSymbolForParsing
	ErrEvaluatorTimestampFormatPatternDuplicateFields
	ErrEvaluatorTimestampFormatPatternHourClockAmPmMismatch
	ErrEvaluatorUnterminatedTimestampFormatPatternToken
	ErrEvaluatorInvalidTimestampFormatPatternToken
	ErrEvaluatorInvalidTimestampFormatPatternSymbol
	ErrEvaluatorBindingDoesNotExist
	ErrMissingHeaders
	ErrInvalidColumnIndex

	ErrAdminConfigNotificationTargetsFailed
	ErrAdminProfilerNotEnabled
	ErrInvalidDecompressedSize
	ErrAddUserInvalidArgument
	ErrAdminAccountNotEligible
	ErrServiceAccountNotFound
	ErrPostPolicyConditionInvalidFormat
)

// error code to Error structure, these fields carry respective
// descriptions for all the error responses.
var errorCodes = errorCodeMap{
	ErrInvalidCopyDest: {
		ErrCode:        ErrInvalidCopyDest,
		Code:           "InvalidRequest",
		Description:    "This copy request is illegal because it is trying to copy an object to itself without changing the object's metadata, storage class, website redirect location or encryption attributes.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidCopySource: {
		ErrCode:        ErrInvalidCopySource,
		Code:           "InvalidArgument",
		Description:    "Copy Source must mention the source bucket and key: sourcebucket/sourcekey.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidMetadataDirective: {
		ErrCode:        ErrInvalidMetadataDirective,
		Code:           "InvalidArgument",
		Description:    "Unknown metadata directive.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidStorageClass: {
		ErrCode:        ErrInvalidStorageClass,
		Code:           "InvalidStorageClass",
		Description:    "Invalid storage class.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidRequestBody: {
		ErrCode:        ErrInvalidRequestBody,
		Code:           "InvalidArgument",
		Description:    "Body shouldn't be set for this request.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidMaxUploads: {
		ErrCode:        ErrInvalidMaxUploads,
		Code:           "InvalidArgument",
		Description:    "Argument max-uploads must be an integer between 0 and 2147483647",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidMaxKeys: {
		ErrCode:        ErrInvalidMaxKeys,
		Code:           "InvalidArgument",
		Description:    "Argument maxKeys must be an integer between 0 and 2147483647",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidEncodingMethod: {
		ErrCode:        ErrInvalidEncodingMethod,
		Code:           "InvalidArgument",
		Description:    "Invalid Encoding Method specified in Request",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidMaxParts: {
		ErrCode:        ErrInvalidMaxParts,
		Code:           "InvalidArgument",
		Description:    "Argument max-parts must be an integer between 0 and 2147483647",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidPartNumberMarker: {
		ErrCode:        ErrInvalidPartNumberMarker,
		Code:           "InvalidArgument",
		Description:    "Argument partNumberMarker must be an integer.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidPolicyDocument: {
		ErrCode:        ErrInvalidPolicyDocument,
		Code:           "InvalidPolicyDocument",
		Description:    "The content of the form does not meet the conditions specified in the policy document.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAccessDenied: {
		ErrCode:        ErrAccessDenied,
		Code:           "AccessDenied",
		Description:    "Access Denied.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrBadDigest: {
		ErrCode:        ErrBadDigest,
		Code:           "BadDigest",
		Description:    "The Content-Md5 you specified did not match what we received.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEntityTooSmall: {
		ErrCode:        ErrEntityTooSmall,
		Code:           "EntityTooSmall",
		Description:    "Your proposed upload is smaller than the minimum allowed object size.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEntityTooLarge: {
		ErrCode:        ErrEntityTooLarge,
		Code:           "EntityTooLarge",
		Description:    "Your proposed upload exceeds the maximum allowed object size.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrPolicyTooLarge: {
		ErrCode:        ErrPolicyTooLarge,
		Code:           "PolicyTooLarge",
		Description:    "Policy exceeds the maximum allowed document size.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrIllegalVersioningConfigurationException: {
		ErrCode:        ErrIllegalVersioningConfigurationException,
		Code:           "IllegalVersioningConfigurationException",
		Description:    "Indicates that the versioning configuration specified in the request is invalid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrIncompleteBody: {
		ErrCode:        ErrIncompleteBody,
		Code:           "IncompleteBody",
		Description:    "You did not provide the number of bytes specified by the Content-Length HTTP header.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInternalError: {
		ErrCode:        ErrInternalError,
		Code:           "InternalError",
		Description:    "We encountered an internal error, please try again.",
		HTTPStatusCode: http.StatusInternalServerError,
	},
	ErrInvalidAccessKeyID: {
		ErrCode:        ErrInvalidAccessKeyID,
		Code:           "InvalidAccessKeyId",
		Description:    "The Access Key Id you provided does not exist in our records.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrInvalidBucketName: {
		ErrCode:        ErrInvalidBucketName,
		Code:           "InvalidBucketName",
		Description:    "The specified bucket is not valid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidDigest: {
		ErrCode:        ErrInvalidDigest,
		Code:           "InvalidDigest",
		Description:    "The Content-Md5 you specified is not valid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidRange: {
		ErrCode:        ErrInvalidRange,
		Code:           "InvalidRange",
		Description:    "The requested range is not satisfiable",
		HTTPStatusCode: http.StatusRequestedRangeNotSatisfiable,
	},
	ErrMalformedXML: {
		ErrCode:        ErrMalformedXML,
		Code:           "MalformedXML",
		Description:    "The XML you provided was not well-formed or did not validate against our published schema.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingContentLength: {
		ErrCode:        ErrMissingContentLength,
		Code:           "MissingContentLength",
		Description:    "You must provide the Content-Length HTTP header.",
		HTTPStatusCode: http.StatusLengthRequired,
	},
	ErrMissingContentMD5: {
		ErrCode:        ErrMissingContentMD5,
		Code:           "MissingContentMD5",
		Description:    "Missing required header for this request: Content-Md5.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingSecurityHeader: {
		ErrCode:        ErrMissingSecurityHeader,
		Code:           "MissingSecurityHeader",
		Description:    "Your request was missing a required header",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingRequestBodyError: {
		ErrCode:        ErrMissingRequestBodyError,
		Code:           "MissingRequestBodyError",
		Description:    "Request body is empty.",
		HTTPStatusCode: http.StatusLengthRequired,
	},
	ErrNoSuchBucket: {
		ErrCode:        ErrNoSuchBucket,
		Code:           "NoSuchBucket",
		Description:    "The specified bucket does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchBucketPolicy: {
		ErrCode:        ErrNoSuchBucketPolicy,
		Code:           "NoSuchBucketPolicy",
		Description:    "The bucket policy does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchBucketLifecycle: {
		ErrCode:        ErrNoSuchBucketLifecycle,
		Code:           "NoSuchBucketLifecycle",
		Description:    "The bucket lifecycle configuration does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchLifecycleConfiguration: {
		ErrCode:        ErrNoSuchLifecycleConfiguration,
		Code:           "NoSuchLifecycleConfiguration",
		Description:    "The lifecycle configuration does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchBucketSSEConfig: {
		ErrCode:        ErrNoSuchBucketSSEConfig,
		Code:           "ServerSideEncryptionConfigurationNotFoundError",
		Description:    "The server side encryption configuration was not found",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchKey: {
		ErrCode:        ErrNoSuchKey,
		Code:           "NoSuchKey",
		Description:    "The specified key does not exist.",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchUpload: {
		ErrCode:        ErrNoSuchUpload,
		Code:           "NoSuchUpload",
		Description:    "The specified multipart upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchVersion: {
		ErrCode:        ErrNoSuchVersion,
		Code:           "NoSuchVersion",
		Description:    "Indicates that the version ID specified in the request does not match an existing version.",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrInvalidVersion: {
		ErrCode:        ErrInvalidVersion,
		Code:           "InvalidArgument",
		Description:    "Invalid version id specified",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidArgument: {
		ErrCode:        ErrInvalidArgument,
		Code:           "InvalidArgument",
		Description:    "The specified argument was invalid",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidTagKey: {
		ErrCode:        ErrInvalidTagKey,
		Code:           "InvalidTag",
		Description:    "The TagValue you have provided is invalid",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidTagValue: {
		ErrCode:        ErrInvalidTagValue,
		Code:           "InvalidTag",
		Description:    "The TagKey you have provided is invalid",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidTagsSizeExceed: {
		ErrCode:        ErrInvalidTagsSizeExceed,
		Code:           "BadRequest",
		Description:    "Object tags cannot be greater than 10",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrNotImplemented: {
		ErrCode:        ErrNotImplemented,
		Code:           "NotImplemented",
		Description:    "A header you provided implies functionality that is not implemented",
		HTTPStatusCode: http.StatusNotImplemented,
	},
	ErrPreconditionFailed: {
		ErrCode:        ErrPreconditionFailed,
		Code:           "PreconditionFailed",
		Description:    "At least one of the pre-conditions you specified did not hold",
		HTTPStatusCode: http.StatusPreconditionFailed,
	},
	ErrNotModified: {
		ErrCode:        ErrNotModified,
		Code:           "NotModified",
		Description:    "The resource was not changed.",
		HTTPStatusCode: http.StatusNotModified,
	},
	ErrRequestTimeTooSkewed: {
		ErrCode:        ErrRequestTimeTooSkewed,
		Code:           "RequestTimeTooSkewed",
		Description:    "The difference between the request time and the server's time is too large.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrSignatureDoesNotMatch: {
		ErrCode:        ErrSignatureDoesNotMatch,
		Code:           "SignatureDoesNotMatch",
		Description:    "The request signature we calculated does not match the signature you provided. Check your key and signing method.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrMethodNotAllowed: {
		ErrCode:        ErrMethodNotAllowed,
		Code:           "MethodNotAllowed",
		Description:    "The specified method is not allowed against this resource.",
		HTTPStatusCode: http.StatusMethodNotAllowed,
	},
	ErrInvalidPart: {
		ErrCode:        ErrInvalidPart,
		Code:           "InvalidPart",
		Description:    "One or more of the specified parts could not be found.  The part may not have been uploaded, or the specified entity tag may not match the part's entity tag.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidPartOrder: {
		ErrCode:        ErrInvalidPartOrder,
		Code:           "InvalidPartOrder",
		Description:    "The list of parts was not in ascending order. The parts list must be specified in order by part number.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidObjectState: {
		ErrCode:        ErrInvalidObjectState,
		Code:           "InvalidObjectState",
		Description:    "The operation is not valid for the current state of the object.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrAuthorizationHeaderMalformed: {
		ErrCode:        ErrAuthorizationHeaderMalformed,
		Code:           "AuthorizationHeaderMalformed",
		Description:    "The authorization header is malformed; the region is wrong; expecting 'us-east-1'.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMalformedPOSTRequest: {
		ErrCode:        ErrMalformedPOSTRequest,
		Code:           "MalformedPOSTRequest",
		Description:    "The body of your POST request is not well-formed multipart/form-data.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrPOSTFileRequired: {
		ErrCode:        ErrPOSTFileRequired,
		Code:           "InvalidArgument",
		Description:    "POST requires exactly one file upload per request.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSignatureVersionNotSupported: {
		ErrCode:        ErrSignatureVersionNotSupported,
		Code:           "InvalidRequest",
		Description:    "The authorization mechanism you have provided is not supported. Please use AWS4-HMAC-SHA256.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrBucketNotEmpty: {
		ErrCode:        ErrBucketNotEmpty,
		Code:           "BucketNotEmpty",
		Description:    "The bucket you tried to delete is not empty",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrBucketAlreadyExists: {
		ErrCode:        ErrBucketAlreadyExists,
		Code:           "BucketAlreadyExists",
		Description:    "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again.",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrAllAccessDisabled: {
		ErrCode:        ErrAllAccessDisabled,
		Code:           "AllAccessDisabled",
		Description:    "All access to this bucket has been disabled.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrMalformedPolicy: {
		ErrCode:        ErrMalformedPolicy,
		Code:           "MalformedPolicy",
		Description:    "Policy has invalid resource.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingFields: {
		ErrCode:        ErrMissingFields,
		Code:           "MissingFields",
		Description:    "Missing fields in request.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingCredTag: {
		ErrCode:        ErrMissingCredTag,
		Code:           "InvalidRequest",
		Description:    "Missing Credential field for this request.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCredMalformed: {
		ErrCode:        ErrCredMalformed,
		Code:           "AuthorizationQueryParametersError",
		Description:    "Error parsing the X-Amz-Credential parameter; the Credential is mal-formed; expecting \"<YOUR-AKID>/YYYYMMDD/REGION/SERVICE/aws4_request\".",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMalformedDate: {
		ErrCode:        ErrMalformedDate,
		Code:           "MalformedDate",
		Description:    "Invalid date format header, expected to be in ISO8601, RFC1123 or RFC1123Z time format.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMalformedPresignedDate: {
		ErrCode:        ErrMalformedPresignedDate,
		Code:           "AuthorizationQueryParametersError",
		Description:    "X-Amz-Date must be in the ISO8601 Long Format \"yyyyMMdd'T'HHmmss'Z'\"",
		HTTPStatusCode: http.StatusBadRequest,
	},
	// FIXME: Should contain the invalid param set as seen in https://github.com/minio/minio/issues/2385.
	// right Description:    "Error parsing the X-Amz-Credential parameter; incorrect date format \"%s\". This date in the credential must be in the format \"yyyyMMdd\".",
	// Need changes to make sure variable messages can be constructed.
	ErrMalformedCredentialDate: {
		ErrCode:        ErrMalformedCredentialDate,
		Code:           "AuthorizationQueryParametersError",
		Description:    "Error parsing the X-Amz-Credential parameter; incorrect date format \"%s\". This date in the credential must be in the format \"yyyyMMdd\".",
		HTTPStatusCode: http.StatusBadRequest,
	},
	// FIXME: Should contain the invalid param set as seen in https://github.com/minio/minio/issues/2385.
	// right Description:    "Error parsing the X-Amz-Credential parameter; the region 'us-east-' is wrong; expecting 'us-east-1'".
	// Need changes to make sure variable messages can be constructed.
	ErrMalformedCredentialRegion: {
		ErrCode:        ErrMalformedCredentialRegion,
		Code:           "AuthorizationQueryParametersError",
		Description:    "Error parsing the X-Amz-Credential parameter; the region is wrong;",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidRegion: {
		ErrCode:        ErrInvalidRegion,
		Code:           "InvalidRegion",
		Description:    "Region does not match.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	// FIXME: Should contain the invalid param set as seen in https://github.com/minio/minio/issues/2385.
	// right Description:   "Error parsing the X-Amz-Credential parameter; incorrect service \"s4\". This endpoint belongs to \"s3\".".
	// Need changes to make sure variable messages can be constructed.
	ErrInvalidServiceS3: {
		ErrCode:        ErrInvalidServiceS3,
		Code:           "AuthorizationParametersError",
		Description:    "Error parsing the Credential/X-Amz-Credential parameter; incorrect service. This endpoint belongs to \"s3\".",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidServiceSTS: {
		ErrCode:        ErrInvalidServiceSTS,
		Code:           "AuthorizationParametersError",
		Description:    "Error parsing the Credential parameter; incorrect service. This endpoint belongs to \"sts\".",
		HTTPStatusCode: http.StatusBadRequest,
	},
	// FIXME: Should contain the invalid param set as seen in https://github.com/minio/minio/issues/2385.
	// Description:   "Error parsing the X-Amz-Credential parameter; incorrect terminal "aws4_reque". This endpoint uses "aws4_request".
	// Need changes to make sure variable messages can be constructed.
	ErrInvalidRequestVersion: {
		ErrCode:        ErrInvalidRequestVersion,
		Code:           "AuthorizationQueryParametersError",
		Description:    "Error parsing the X-Amz-Credential parameter; incorrect terminal. This endpoint uses \"aws4_request\".",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingSignTag: {
		ErrCode:        ErrMissingSignTag,
		Code:           "AccessDenied",
		Description:    "Signature header missing Signature field.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingSignHeadersTag: {
		ErrCode:        ErrMissingSignHeadersTag,
		Code:           "InvalidArgument",
		Description:    "Signature header missing SignedHeaders field.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMalformedExpires: {
		ErrCode:        ErrMalformedExpires,
		Code:           "AuthorizationQueryParametersError",
		Description:    "X-Amz-Expires should be a number",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrNegativeExpires: {
		ErrCode:        ErrNegativeExpires,
		Code:           "AuthorizationQueryParametersError",
		Description:    "X-Amz-Expires must be non-negative",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAuthHeaderEmpty: {
		ErrCode:        ErrAuthHeaderEmpty,
		Code:           "InvalidArgument",
		Description:    "Authorization header is invalid -- one and only one ' ' (space) required.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingDateHeader: {
		ErrCode:        ErrMissingDateHeader,
		Code:           "AccessDenied",
		Description:    "AWS authentication requires a valid Date or x-amz-date header",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidQuerySignatureAlgo: {
		ErrCode:        ErrInvalidQuerySignatureAlgo,
		Code:           "AuthorizationQueryParametersError",
		Description:    "X-Amz-Algorithm only supports \"AWS4-HMAC-SHA256\".",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrExpiredPresignRequest: {
		ErrCode:        ErrExpiredPresignRequest,
		Code:           "AccessDenied",
		Description:    "Request has expired",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrRequestNotReadyYet: {
		ErrCode:        ErrRequestNotReadyYet,
		Code:           "AccessDenied",
		Description:    "Request is not valid yet",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrSlowDown: {
		ErrCode:        ErrSlowDown,
		Code:           "SlowDown",
		Description:    "Please reduce your request",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrInvalidPrefixMarker: {
		ErrCode:        ErrInvalidPrefixMarker,
		Code:           "InvalidPrefixMarker",
		Description:    "Invalid marker prefix combination",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrBadRequest: {
		ErrCode:        ErrBadRequest,
		Code:           "BadRequest",
		Description:    "400 BadRequest",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrKeyTooLongError: {
		ErrCode:        ErrKeyTooLongError,
		Code:           "KeyTooLongError",
		Description:    "Your key is too long",
		HTTPStatusCode: http.StatusBadRequest,
	},

	// FIXME: Actual XML error response also contains the header which missed in list of signed header parameters.
	ErrUnsignedHeaders: {
		ErrCode:        ErrUnsignedHeaders,
		Code:           "AccessDenied",
		Description:    "There were headers present in the request which were not signed",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidQueryParams: {
		ErrCode:        ErrInvalidQueryParams,
		Code:           "AuthorizationQueryParametersError",
		Description:    "Query-string authentication version 4 requires the X-Amz-Algorithm, X-Amz-Credential, X-Amz-Signature, X-Amz-Date, X-Amz-SignedHeaders, and X-Amz-Expires parameters.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrBucketAlreadyOwnedByYou: {
		ErrCode:        ErrBucketAlreadyOwnedByYou,
		Code:           "BucketAlreadyOwnedByYou",
		Description:    "Your previous request to create the named bucket succeeded and you already own it.",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrInvalidDuration: {
		ErrCode:        ErrInvalidDuration,
		Code:           "InvalidDuration",
		Description:    "Duration provided in the request is invalid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidBucketObjectLockConfiguration: {
		ErrCode:        ErrInvalidBucketObjectLockConfiguration,
		Code:           "InvalidRequest",
		Description:    "Bucket is missing ObjectLockConfiguration",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrBucketTaggingNotFound: {
		ErrCode:        ErrBucketTaggingNotFound,
		Code:           "NoSuchTagSet",
		Description:    "The TagSet does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrObjectLockConfigurationNotFound: {
		ErrCode:        ErrObjectLockConfigurationNotFound,
		Code:           "ObjectLockConfigurationNotFoundError",
		Description:    "Object Lock configuration does not exist for this bucket",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrObjectLockConfigurationNotAllowed: {
		ErrCode:        ErrObjectLockConfigurationNotAllowed,
		Code:           "InvalidBucketState",
		Description:    "Object Lock configuration cannot be enabled on existing buckets",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrNoSuchCORSConfiguration: {
		ErrCode:        ErrNoSuchCORSConfiguration,
		Code:           "NoSuchCORSConfiguration",
		Description:    "The CORS configuration does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchWebsiteConfiguration: {
		ErrCode:        ErrNoSuchWebsiteConfiguration,
		Code:           "NoSuchWebsiteConfiguration",
		Description:    "The specified bucket does not have a website configuration",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrReplicationConfigurationNotFoundError: {
		ErrCode:        ErrReplicationConfigurationNotFoundError,
		Code:           "ReplicationConfigurationNotFoundError",
		Description:    "The replication configuration was not found",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchObjectLockConfiguration: {
		ErrCode:        ErrNoSuchObjectLockConfiguration,
		Code:           "NoSuchObjectLockConfiguration",
		Description:    "The specified object does not have a ObjectLock configuration",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrObjectLocked: {
		ErrCode:        ErrObjectLocked,
		Code:           "InvalidRequest",
		Description:    "Object is WORM protected and cannot be overwritten",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidRetentionDate: {
		ErrCode:        ErrInvalidRetentionDate,
		Code:           "InvalidRequest",
		Description:    "Date must be provided in ISO 8601 format",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrPastObjectLockRetainDate: {
		ErrCode:        ErrPastObjectLockRetainDate,
		Code:           "InvalidRequest",
		Description:    "the retain until date must be in the future",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrUnknownWORMModeDirective: {
		ErrCode:        ErrUnknownWORMModeDirective,
		Code:           "InvalidRequest",
		Description:    "unknown wormMode directive",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrObjectLockInvalidHeaders: {
		ErrCode:        ErrObjectLockInvalidHeaders,
		Code:           "InvalidRequest",
		Description:    "x-amz-object-lock-retain-until-date and x-amz-object-lock-mode must both be supplied",
		HTTPStatusCode: http.StatusBadRequest,
	},
	// Bucket notification related errors.
	ErrEventNotification: {
		ErrCode:        ErrEventNotification,
		Code:           "InvalidArgument",
		Description:    "A specified event is not supported for notifications.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrARNNotification: {
		ErrCode:        ErrARNNotification,
		Code:           "InvalidArgument",
		Description:    "A specified destination ARN does not exist or is not well-formed. Verify the destination ARN.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrRegionNotification: {
		ErrCode:        ErrRegionNotification,
		Code:           "InvalidArgument",
		Description:    "A specified destination is in a different region than the bucket. You must use a destination that resides in the same region as the bucket.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrOverlappingFilterNotification: {
		ErrCode:        ErrOverlappingFilterNotification,
		Code:           "InvalidArgument",
		Description:    "An object key name filtering rule defined with overlapping prefixes, overlapping suffixes, or overlapping combinations of prefixes and suffixes for the same event types.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrFilterNameInvalid: {
		ErrCode:        ErrFilterNameInvalid,
		Code:           "InvalidArgument",
		Description:    "filter rule name must be either prefix or suffix",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrFilterNamePrefix: {
		ErrCode:        ErrFilterNamePrefix,
		Code:           "InvalidArgument",
		Description:    "Cannot specify more than one prefix rule in a filter.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrFilterNameSuffix: {
		ErrCode:        ErrFilterNameSuffix,
		Code:           "InvalidArgument",
		Description:    "Cannot specify more than one suffix rule in a filter.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrFilterValueInvalid: {
		ErrCode:        ErrFilterValueInvalid,
		Code:           "InvalidArgument",
		Description:    "Size of filter rule value cannot exceed 1024 bytes in UTF-8 representation",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrOverlappingConfigs: {
		ErrCode:        ErrOverlappingConfigs,
		Code:           "InvalidArgument",
		Description:    "Configurations overlap. Configurations on the same bucket cannot share a common event type.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrUnsupportedNotification: {
		ErrCode:        ErrUnsupportedNotification,
		Code:           "UnsupportedNotification",
		Description:    "MinIO server does not support Topic or Cloud Function based notifications.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidCopyPartRange: {
		ErrCode:        ErrInvalidCopyPartRange,
		Code:           "InvalidArgument",
		Description:    "The x-amz-copy-source-range value must be of the form bytes=first-last where first and last are the zero-based offsets of the first and last bytes to copy",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidCopyPartRangeSource: {
		ErrCode:        ErrInvalidCopyPartRangeSource,
		Code:           "InvalidArgument",
		Description:    "Range specified is not valid for source object",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMetadataTooLarge: {
		ErrCode:        ErrMetadataTooLarge,
		Code:           "InvalidArgument",
		Description:    "Your metadata headers exceed the maximum allowed metadata size.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidTagDirective: {
		ErrCode:        ErrInvalidTagDirective,
		Code:           "InvalidArgument",
		Description:    "Unknown tag directive.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrNotSupported: {
		ErrCode:        ErrNotSupported,
		Code:           "BadRequest",
		Description:    "Not supported by NeoFS S3 Gateway",
		HTTPStatusCode: http.StatusNotImplemented,
	},
	ErrInvalidEncryptionMethod: {
		ErrCode:        ErrInvalidEncryptionMethod,
		Code:           "InvalidRequest",
		Description:    "The encryption method specified is not supported",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInsecureSSECustomerRequest: {
		ErrCode:        ErrInsecureSSECustomerRequest,
		Code:           "InvalidRequest",
		Description:    "Requests specifying Server Side Encryption with Customer provided keys must be made over a secure connection.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSSEMultipartEncrypted: {
		ErrCode:        ErrSSEMultipartEncrypted,
		Code:           "InvalidRequest",
		Description:    "The multipart upload initiate requested encryption. Subsequent part requests must include the appropriate encryption parameters.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSSEEncryptedObject: {
		ErrCode:        ErrSSEEncryptedObject,
		Code:           "InvalidRequest",
		Description:    "The object was stored using a form of Server Side Encryption. The correct parameters must be provided to retrieve the object.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidEncryptionParameters: {
		ErrCode:        ErrInvalidEncryptionParameters,
		Code:           "InvalidRequest",
		Description:    "The encryption parameters are not applicable to this object.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidSSECustomerAlgorithm: {
		ErrCode:        ErrInvalidSSECustomerAlgorithm,
		Code:           "InvalidArgument",
		Description:    "Requests specifying Server Side Encryption with Customer provided keys must provide a valid encryption algorithm.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidSSECustomerKey: {
		ErrCode:        ErrInvalidSSECustomerKey,
		Code:           "InvalidArgument",
		Description:    "The secret key was invalid for the specified algorithm.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingSSECustomerKey: {
		ErrCode:        ErrMissingSSECustomerKey,
		Code:           "InvalidArgument",
		Description:    "Requests specifying Server Side Encryption with Customer provided keys must provide an appropriate secret key.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingSSECustomerKeyMD5: {
		ErrCode:        ErrMissingSSECustomerKeyMD5,
		Code:           "InvalidArgument",
		Description:    "Requests specifying Server Side Encryption with Customer provided keys must provide the client calculated MD5 of the secret key.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSSECustomerKeyMD5Mismatch: {
		ErrCode:        ErrSSECustomerKeyMD5Mismatch,
		Code:           "InvalidArgument",
		Description:    "The calculated MD5 hash of the key did not match the hash that was provided.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidSSECustomerParameters: {
		ErrCode:        ErrInvalidSSECustomerParameters,
		Code:           "InvalidArgument",
		Description:    "The provided encryption parameters did not match the ones used originally.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrIncompatibleEncryptionMethod: {
		ErrCode:        ErrIncompatibleEncryptionMethod,
		Code:           "InvalidArgument",
		Description:    "Server side encryption specified with both SSE-C and SSE-S3 headers",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrKMSNotConfigured: {
		ErrCode:        ErrKMSNotConfigured,
		Code:           "InvalidArgument",
		Description:    "Server side encryption specified but KMS is not configured",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrKMSAuthFailure: {
		ErrCode:        ErrKMSAuthFailure,
		Code:           "InvalidArgument",
		Description:    "Server side encryption specified but KMS authorization failed",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrNoAccessKey: {
		ErrCode:        ErrNoAccessKey,
		Code:           "AccessDenied",
		Description:    "No AWSAccessKey was presented",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrInvalidToken: {
		ErrCode:        ErrInvalidToken,
		Code:           "InvalidTokenId",
		Description:    "The security token included in the request is invalid",
		HTTPStatusCode: http.StatusForbidden,
	},

	// S3 extensions.
	ErrContentSHA256Mismatch: {
		ErrCode:        ErrContentSHA256Mismatch,
		Code:           "XAmzContentSHA256Mismatch",
		Description:    "The provided 'x-amz-content-sha256' header does not match what was computed.",
		HTTPStatusCode: http.StatusBadRequest,
	},

	// MinIO extensions.
	ErrStorageFull: {
		ErrCode:        ErrStorageFull,
		Code:           "XMinioStorageFull",
		Description:    "Storage backend has reached its minimum free disk threshold. Please delete a few objects to proceed.",
		HTTPStatusCode: http.StatusInsufficientStorage,
	},
	ErrParentIsObject: {
		ErrCode:        ErrParentIsObject,
		Code:           "XMinioParentIsObject",
		Description:    "Object-prefix is already an object, please choose a different object-prefix name.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrRequestBodyParse: {
		ErrCode:        ErrRequestBodyParse,
		Code:           "XMinioRequestBodyParse",
		Description:    "The request body failed to parse.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrObjectExistsAsDirectory: {
		ErrCode:        ErrObjectExistsAsDirectory,
		Code:           "XMinioObjectExistsAsDirectory",
		Description:    "Object name already exists as a directory.",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrInvalidObjectName: {
		ErrCode:        ErrInvalidObjectName,
		Code:           "XMinioInvalidObjectName",
		Description:    "Object name contains unsupported characters.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidObjectNamePrefixSlash: {
		ErrCode:        ErrInvalidObjectNamePrefixSlash,
		Code:           "XMinioInvalidObjectName",
		Description:    "Object name contains a leading slash.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidResourceName: {
		ErrCode:        ErrInvalidResourceName,
		Code:           "XMinioInvalidResourceName",
		Description:    "Resource name contains bad components such as \"..\" or \".\".",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrServerNotInitialized: {
		ErrCode:        ErrServerNotInitialized,
		Code:           "XMinioServerNotInitialized",
		Description:    "Server not initialized, please try again.",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrMalformedJSON: {
		ErrCode:        ErrMalformedJSON,
		Code:           "XMinioMalformedJSON",
		Description:    "The JSON you provided was not well-formed or did not validate against our published format.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAdminNoSuchUser: {
		ErrCode:        ErrAdminNoSuchUser,
		Code:           "XMinioAdminNoSuchUser",
		Description:    "The specified user does not exist.",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrAdminNoSuchGroup: {
		ErrCode:        ErrAdminNoSuchGroup,
		Code:           "XMinioAdminNoSuchGroup",
		Description:    "The specified group does not exist.",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrAdminGroupNotEmpty: {
		ErrCode:        ErrAdminGroupNotEmpty,
		Code:           "XMinioAdminGroupNotEmpty",
		Description:    "The specified group is not empty - cannot remove it.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAdminNoSuchPolicy: {
		ErrCode:        ErrAdminNoSuchPolicy,
		Code:           "XMinioAdminNoSuchPolicy",
		Description:    "The canned policy does not exist.",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrAdminInvalidArgument: {
		ErrCode:        ErrAdminInvalidArgument,
		Code:           "XMinioAdminInvalidArgument",
		Description:    "Invalid arguments specified.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAdminInvalidAccessKey: {
		ErrCode:        ErrAdminInvalidAccessKey,
		Code:           "XMinioAdminInvalidAccessKey",
		Description:    "The access key is invalid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAdminInvalidSecretKey: {
		ErrCode:        ErrAdminInvalidSecretKey,
		Code:           "XMinioAdminInvalidSecretKey",
		Description:    "The secret key is invalid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAdminConfigNoQuorum: {
		ErrCode:        ErrAdminConfigNoQuorum,
		Code:           "XMinioAdminConfigNoQuorum",
		Description:    "Configuration update failed because server quorum was not met",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrAdminConfigTooLarge: {
		ErrCode: ErrAdminConfigTooLarge,
		Code:    "XMinioAdminConfigTooLarge",
		Description: fmt.Sprintf("Configuration data provided exceeds the allowed maximum of %d bytes",
			maxEConfigJSONSize),
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAdminConfigBadJSON: {
		ErrCode:        ErrAdminConfigBadJSON,
		Code:           "XMinioAdminConfigBadJSON",
		Description:    "JSON configuration provided is of incorrect format",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAdminConfigDuplicateKeys: {
		ErrCode:        ErrAdminConfigDuplicateKeys,
		Code:           "XMinioAdminConfigDuplicateKeys",
		Description:    "JSON configuration provided has objects with duplicate keys",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAdminConfigNotificationTargetsFailed: {
		ErrCode:        ErrAdminConfigNotificationTargetsFailed,
		Code:           "XMinioAdminNotificationTargetsTestFailed",
		Description:    "Configuration update failed due an unsuccessful attempt to connect to one or more notification servers",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAdminProfilerNotEnabled: {
		ErrCode:        ErrAdminProfilerNotEnabled,
		Code:           "XMinioAdminProfilerNotEnabled",
		Description:    "Unable to perform the requested operation because profiling is not enabled",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAdminCredentialsMismatch: {
		ErrCode:        ErrAdminCredentialsMismatch,
		Code:           "XMinioAdminCredentialsMismatch",
		Description:    "Credentials in config mismatch with server environment variables",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrAdminBucketQuotaExceeded: {
		ErrCode:        ErrAdminBucketQuotaExceeded,
		Code:           "XMinioAdminBucketQuotaExceeded",
		Description:    "Bucket quota exceeded",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAdminNoSuchQuotaConfiguration: {
		ErrCode:        ErrAdminNoSuchQuotaConfiguration,
		Code:           "XMinioAdminNoSuchQuotaConfiguration",
		Description:    "The quota configuration does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrAdminBucketQuotaDisabled: {
		ErrCode:        ErrAdminBucketQuotaDisabled,
		Code:           "XMinioAdminBucketQuotaDisabled",
		Description:    "Quota specified but disk usage crawl is disabled on MinIO server",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInsecureClientRequest: {
		ErrCode:        ErrInsecureClientRequest,
		Code:           "XMinioInsecureClientRequest",
		Description:    "Cannot respond to plain-text request from TLS-encrypted server",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrOperationTimedOut: {
		ErrCode:        ErrOperationTimedOut,
		Code:           "RequestTimeout",
		Description:    "A timeout occurred while trying to lock a resource, please reduce your request rate",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrOperationMaxedOut: {
		ErrCode:        ErrOperationMaxedOut,
		Code:           "SlowDown",
		Description:    "A timeout exceeded while waiting to proceed with the request, please reduce your request rate",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrUnsupportedMetadata: {
		ErrCode:        ErrUnsupportedMetadata,
		Code:           "InvalidArgument",
		Description:    "Your metadata headers are not supported.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrObjectTampered: {
		ErrCode:        ErrObjectTampered,
		Code:           "XMinioObjectTampered",
		Description:    "The requested object was modified and may be compromised",
		HTTPStatusCode: http.StatusPartialContent,
	},
	ErrMaximumExpires: {
		ErrCode:        ErrMaximumExpires,
		Code:           "AuthorizationQueryParametersError",
		Description:    "X-Amz-Expires must be less than a week (in seconds); that is, the given X-Amz-Expires must be less than 604800 seconds",
		HTTPStatusCode: http.StatusBadRequest,
	},

	// Generic Invalid-Request error. Should be used for response errors only for unlikely
	// corner case errors for which introducing new ErrorCode is not worth it. LogIf()
	// should be used to log the error at the source of the error for debugging purposes.
	ErrInvalidRequest: {
		ErrCode:        ErrInvalidRequest,
		Code:           "InvalidRequest",
		Description:    "Invalid Request",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrHealNotImplemented: {
		ErrCode:        ErrHealNotImplemented,
		Code:           "XMinioHealNotImplemented",
		Description:    "This server does not implement heal functionality.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrHealNoSuchProcess: {
		ErrCode:        ErrHealNoSuchProcess,
		Code:           "XMinioHealNoSuchProcess",
		Description:    "No such heal process is running on the server",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrHealInvalidClientToken: {
		ErrCode:        ErrHealInvalidClientToken,
		Code:           "XMinioHealInvalidClientToken",
		Description:    "Client token mismatch",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrHealMissingBucket: {
		ErrCode:        ErrHealMissingBucket,
		Code:           "XMinioHealMissingBucket",
		Description:    "A heal start request with a non-empty object-prefix parameter requires a bucket to be specified.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrHealAlreadyRunning: {
		ErrCode:        ErrHealAlreadyRunning,
		Code:           "XMinioHealAlreadyRunning",
		Description:    "",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrHealOverlappingPaths: {
		ErrCode:        ErrHealOverlappingPaths,
		Code:           "XMinioHealOverlappingPaths",
		Description:    "",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrBackendDown: {
		ErrCode:        ErrBackendDown,
		Code:           "XMinioBackendDown",
		Description:    "Object storage backend is unreachable",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrIncorrectContinuationToken: {
		ErrCode:        ErrIncorrectContinuationToken,
		Code:           "InvalidArgument",
		Description:    "The continuation token provided is incorrect",
		HTTPStatusCode: http.StatusBadRequest,
	},

	// S3 Select API Errors
	ErrEmptyRequestBody: {
		ErrCode:        ErrEmptyRequestBody,
		Code:           "EmptyRequestBody",
		Description:    "Request body cannot be empty.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrUnsupportedFunction: {
		ErrCode:        ErrUnsupportedFunction,
		Code:           "UnsupportedFunction",
		Description:    "Encountered an unsupported SQL function.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidDataSource: {
		ErrCode:        ErrInvalidDataSource,
		Code:           "InvalidDataSource",
		Description:    "Invalid data source type. Only CSV and JSON are supported at this time.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidExpressionType: {
		ErrCode:        ErrInvalidExpressionType,
		Code:           "InvalidExpressionType",
		Description:    "The ExpressionType is invalid. Only SQL expressions are supported at this time.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrBusy: {
		ErrCode:        ErrBusy,
		Code:           "Busy",
		Description:    "The service is unavailable. Please retry.",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrUnauthorizedAccess: {
		ErrCode:        ErrUnauthorizedAccess,
		Code:           "UnauthorizedAccess",
		Description:    "You are not authorized to perform this operation",
		HTTPStatusCode: http.StatusUnauthorized,
	},
	ErrExpressionTooLong: {
		ErrCode:        ErrExpressionTooLong,
		Code:           "ExpressionTooLong",
		Description:    "The SQL expression is too long: The maximum byte-length for the SQL expression is 256 KB.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrIllegalSQLFunctionArgument: {
		ErrCode:        ErrIllegalSQLFunctionArgument,
		Code:           "IllegalSqlFunctionArgument",
		Description:    "Illegal argument was used in the SQL function.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidKeyPath: {
		ErrCode:        ErrInvalidKeyPath,
		Code:           "InvalidKeyPath",
		Description:    "Key path in the SQL expression is invalid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidCompressionFormat: {
		ErrCode:        ErrInvalidCompressionFormat,
		Code:           "InvalidCompressionFormat",
		Description:    "The file is not in a supported compression format. Only GZIP is supported at this time.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidFileHeaderInfo: {
		ErrCode:        ErrInvalidFileHeaderInfo,
		Code:           "InvalidFileHeaderInfo",
		Description:    "The FileHeaderInfo is invalid. Only NONE, USE, and IGNORE are supported.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidJSONType: {
		ErrCode:        ErrInvalidJSONType,
		Code:           "InvalidJsonType",
		Description:    "The JsonType is invalid. Only DOCUMENT and LINES are supported at this time.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidQuoteFields: {
		ErrCode:        ErrInvalidQuoteFields,
		Code:           "InvalidQuoteFields",
		Description:    "The QuoteFields is invalid. Only ALWAYS and ASNEEDED are supported.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidRequestParameter: {
		ErrCode:        ErrInvalidRequestParameter,
		Code:           "InvalidRequestParameter",
		Description:    "The value of a parameter in SelectRequest element is invalid. Check the service API documentation and try again.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidDataType: {
		ErrCode:        ErrInvalidDataType,
		Code:           "InvalidDataType",
		Description:    "The SQL expression contains an invalid data type.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidTextEncoding: {
		ErrCode:        ErrInvalidTextEncoding,
		Code:           "InvalidTextEncoding",
		Description:    "Invalid encoding type. Only UTF-8 encoding is supported at this time.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidTableAlias: {
		ErrCode:        ErrInvalidTableAlias,
		Code:           "InvalidTableAlias",
		Description:    "The SQL expression contains an invalid table alias.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingRequiredParameter: {
		ErrCode:        ErrMissingRequiredParameter,
		Code:           "MissingRequiredParameter",
		Description:    "The SelectRequest entity is missing a required parameter. Check the service documentation and try again.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrObjectSerializationConflict: {
		ErrCode:        ErrObjectSerializationConflict,
		Code:           "ObjectSerializationConflict",
		Description:    "The SelectRequest entity can only contain one of CSV or JSON. Check the service documentation and try again.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrUnsupportedSQLOperation: {
		ErrCode:        ErrUnsupportedSQLOperation,
		Code:           "UnsupportedSqlOperation",
		Description:    "Encountered an unsupported SQL operation.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrUnsupportedSQLStructure: {
		ErrCode:        ErrUnsupportedSQLStructure,
		Code:           "UnsupportedSqlStructure",
		Description:    "Encountered an unsupported SQL structure. Check the SQL Reference.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrUnsupportedSyntax: {
		ErrCode:        ErrUnsupportedSyntax,
		Code:           "UnsupportedSyntax",
		Description:    "Encountered invalid syntax.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrUnsupportedRangeHeader: {
		ErrCode:        ErrUnsupportedRangeHeader,
		Code:           "UnsupportedRangeHeader",
		Description:    "Range header is not supported for this operation.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrLexerInvalidChar: {
		ErrCode:        ErrLexerInvalidChar,
		Code:           "LexerInvalidChar",
		Description:    "The SQL expression contains an invalid character.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrLexerInvalidOperator: {
		ErrCode:        ErrLexerInvalidOperator,
		Code:           "LexerInvalidOperator",
		Description:    "The SQL expression contains an invalid literal.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrLexerInvalidLiteral: {
		ErrCode:        ErrLexerInvalidLiteral,
		Code:           "LexerInvalidLiteral",
		Description:    "The SQL expression contains an invalid operator.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrLexerInvalidIONLiteral: {
		ErrCode:        ErrLexerInvalidIONLiteral,
		Code:           "LexerInvalidIONLiteral",
		Description:    "The SQL expression contains an invalid operator.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedDatePart: {
		ErrCode:        ErrParseExpectedDatePart,
		Code:           "ParseExpectedDatePart",
		Description:    "Did not find the expected date part in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedKeyword: {
		ErrCode:        ErrParseExpectedKeyword,
		Code:           "ParseExpectedKeyword",
		Description:    "Did not find the expected keyword in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedTokenType: {
		ErrCode:        ErrParseExpectedTokenType,
		Code:           "ParseExpectedTokenType",
		Description:    "Did not find the expected token in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpected2TokenTypes: {
		ErrCode:        ErrParseExpected2TokenTypes,
		Code:           "ParseExpected2TokenTypes",
		Description:    "Did not find the expected token in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedNumber: {
		ErrCode:        ErrParseExpectedNumber,
		Code:           "ParseExpectedNumber",
		Description:    "Did not find the expected number in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedRightParenBuiltinFunctionCall: {
		ErrCode:        ErrParseExpectedRightParenBuiltinFunctionCall,
		Code:           "ParseExpectedRightParenBuiltinFunctionCall",
		Description:    "Did not find the expected right parenthesis character in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedTypeName: {
		ErrCode:        ErrParseExpectedTypeName,
		Code:           "ParseExpectedTypeName",
		Description:    "Did not find the expected type name in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedWhenClause: {
		ErrCode:        ErrParseExpectedWhenClause,
		Code:           "ParseExpectedWhenClause",
		Description:    "Did not find the expected WHEN clause in the SQL expression. CASE is not supported.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnsupportedToken: {
		ErrCode:        ErrParseUnsupportedToken,
		Code:           "ParseUnsupportedToken",
		Description:    "The SQL expression contains an unsupported token.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnsupportedLiteralsGroupBy: {
		ErrCode:        ErrParseUnsupportedLiteralsGroupBy,
		Code:           "ParseUnsupportedLiteralsGroupBy",
		Description:    "The SQL expression contains an unsupported use of GROUP BY.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedMember: {
		ErrCode:        ErrParseExpectedMember,
		Code:           "ParseExpectedMember",
		Description:    "The SQL expression contains an unsupported use of MEMBER.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnsupportedSelect: {
		ErrCode:        ErrParseUnsupportedSelect,
		Code:           "ParseUnsupportedSelect",
		Description:    "The SQL expression contains an unsupported use of SELECT.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnsupportedCase: {
		ErrCode:        ErrParseUnsupportedCase,
		Code:           "ParseUnsupportedCase",
		Description:    "The SQL expression contains an unsupported use of CASE.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnsupportedCaseClause: {
		ErrCode:        ErrParseUnsupportedCaseClause,
		Code:           "ParseUnsupportedCaseClause",
		Description:    "The SQL expression contains an unsupported use of CASE.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnsupportedAlias: {
		ErrCode:        ErrParseUnsupportedAlias,
		Code:           "ParseUnsupportedAlias",
		Description:    "The SQL expression contains an unsupported use of ALIAS.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnsupportedSyntax: {
		ErrCode:        ErrParseUnsupportedSyntax,
		Code:           "ParseUnsupportedSyntax",
		Description:    "The SQL expression contains unsupported syntax.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnknownOperator: {
		ErrCode:        ErrParseUnknownOperator,
		Code:           "ParseUnknownOperator",
		Description:    "The SQL expression contains an invalid operator.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseMissingIdentAfterAt: {
		ErrCode:        ErrParseMissingIdentAfterAt,
		Code:           "ParseMissingIdentAfterAt",
		Description:    "Did not find the expected identifier after the @ symbol in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnexpectedOperator: {
		ErrCode:        ErrParseUnexpectedOperator,
		Code:           "ParseUnexpectedOperator",
		Description:    "The SQL expression contains an unexpected operator.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnexpectedTerm: {
		ErrCode:        ErrParseUnexpectedTerm,
		Code:           "ParseUnexpectedTerm",
		Description:    "The SQL expression contains an unexpected term.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnexpectedToken: {
		ErrCode:        ErrParseUnexpectedToken,
		Code:           "ParseUnexpectedToken",
		Description:    "The SQL expression contains an unexpected token.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnexpectedKeyword: {
		ErrCode:        ErrParseUnexpectedKeyword,
		Code:           "ParseUnexpectedKeyword",
		Description:    "The SQL expression contains an unexpected keyword.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedExpression: {
		ErrCode:        ErrParseExpectedExpression,
		Code:           "ParseExpectedExpression",
		Description:    "Did not find the expected SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedLeftParenAfterCast: {
		ErrCode:        ErrParseExpectedLeftParenAfterCast,
		Code:           "ParseExpectedLeftParenAfterCast",
		Description:    "Did not find expected the left parenthesis in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedLeftParenValueConstructor: {
		ErrCode:        ErrParseExpectedLeftParenValueConstructor,
		Code:           "ParseExpectedLeftParenValueConstructor",
		Description:    "Did not find expected the left parenthesis in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedLeftParenBuiltinFunctionCall: {
		ErrCode:        ErrParseExpectedLeftParenBuiltinFunctionCall,
		Code:           "ParseExpectedLeftParenBuiltinFunctionCall",
		Description:    "Did not find the expected left parenthesis in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedArgumentDelimiter: {
		ErrCode:        ErrParseExpectedArgumentDelimiter,
		Code:           "ParseExpectedArgumentDelimiter",
		Description:    "Did not find the expected argument delimiter in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseCastArity: {
		ErrCode:        ErrParseCastArity,
		Code:           "ParseCastArity",
		Description:    "The SQL expression CAST has incorrect arity.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseInvalidTypeParam: {
		ErrCode:        ErrParseInvalidTypeParam,
		Code:           "ParseInvalidTypeParam",
		Description:    "The SQL expression contains an invalid parameter value.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseEmptySelect: {
		ErrCode:        ErrParseEmptySelect,
		Code:           "ParseEmptySelect",
		Description:    "The SQL expression contains an empty SELECT.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseSelectMissingFrom: {
		ErrCode:        ErrParseSelectMissingFrom,
		Code:           "ParseSelectMissingFrom",
		Description:    "GROUP is not supported in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedIdentForGroupName: {
		ErrCode:        ErrParseExpectedIdentForGroupName,
		Code:           "ParseExpectedIdentForGroupName",
		Description:    "GROUP is not supported in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedIdentForAlias: {
		ErrCode:        ErrParseExpectedIdentForAlias,
		Code:           "ParseExpectedIdentForAlias",
		Description:    "Did not find the expected identifier for the alias in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnsupportedCallWithStar: {
		ErrCode:        ErrParseUnsupportedCallWithStar,
		Code:           "ParseUnsupportedCallWithStar",
		Description:    "Only COUNT with (*) as a parameter is supported in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseNonUnaryAgregateFunctionCall: {
		ErrCode:        ErrParseNonUnaryAgregateFunctionCall,
		Code:           "ParseNonUnaryAgregateFunctionCall",
		Description:    "Only one argument is supported for aggregate functions in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseMalformedJoin: {
		ErrCode:        ErrParseMalformedJoin,
		Code:           "ParseMalformedJoin",
		Description:    "JOIN is not supported in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedIdentForAt: {
		ErrCode:        ErrParseExpectedIdentForAt,
		Code:           "ParseExpectedIdentForAt",
		Description:    "Did not find the expected identifier for AT name in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseAsteriskIsNotAloneInSelectList: {
		ErrCode:        ErrParseAsteriskIsNotAloneInSelectList,
		Code:           "ParseAsteriskIsNotAloneInSelectList",
		Description:    "Other expressions are not allowed in the SELECT list when '*' is used without dot notation in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseCannotMixSqbAndWildcardInSelectList: {
		ErrCode:        ErrParseCannotMixSqbAndWildcardInSelectList,
		Code:           "ParseCannotMixSqbAndWildcardInSelectList",
		Description:    "Cannot mix [] and * in the same expression in a SELECT list in SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseInvalidContextForWildcardInSelectList: {
		ErrCode:        ErrParseInvalidContextForWildcardInSelectList,
		Code:           "ParseInvalidContextForWildcardInSelectList",
		Description:    "Invalid use of * in SELECT list in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrIncorrectSQLFunctionArgumentType: {
		ErrCode:        ErrIncorrectSQLFunctionArgumentType,
		Code:           "IncorrectSqlFunctionArgumentType",
		Description:    "Incorrect type of arguments in function call in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrValueParseFailure: {
		ErrCode:        ErrValueParseFailure,
		Code:           "ValueParseFailure",
		Description:    "Time stamp parse failure in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorInvalidArguments: {
		ErrCode:        ErrEvaluatorInvalidArguments,
		Code:           "EvaluatorInvalidArguments",
		Description:    "Incorrect number of arguments in the function call in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrIntegerOverflow: {
		ErrCode:        ErrIntegerOverflow,
		Code:           "IntegerOverflow",
		Description:    "Int overflow or underflow in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrLikeInvalidInputs: {
		ErrCode:        ErrLikeInvalidInputs,
		Code:           "LikeInvalidInputs",
		Description:    "Invalid argument given to the LIKE clause in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCastFailed: {
		ErrCode:        ErrCastFailed,
		Code:           "CastFailed",
		Description:    "Attempt to convert from one data type to another using CAST failed in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidCast: {
		ErrCode:        ErrInvalidCast,
		Code:           "InvalidCast",
		Description:    "Attempt to convert from one data type to another using CAST failed in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorInvalidTimestampFormatPattern: {
		ErrCode:        ErrEvaluatorInvalidTimestampFormatPattern,
		Code:           "EvaluatorInvalidTimestampFormatPattern",
		Description:    "Time stamp format pattern requires additional fields in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorInvalidTimestampFormatPatternSymbolForParsing: {
		ErrCode:        ErrEvaluatorInvalidTimestampFormatPatternSymbolForParsing,
		Code:           "EvaluatorInvalidTimestampFormatPatternSymbolForParsing",
		Description:    "Time stamp format pattern contains a valid format symbol that cannot be applied to time stamp parsing in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorTimestampFormatPatternDuplicateFields: {
		ErrCode:        ErrEvaluatorTimestampFormatPatternDuplicateFields,
		Code:           "EvaluatorTimestampFormatPatternDuplicateFields",
		Description:    "Time stamp format pattern contains multiple format specifiers representing the time stamp field in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorTimestampFormatPatternHourClockAmPmMismatch: {
		ErrCode:        ErrEvaluatorTimestampFormatPatternHourClockAmPmMismatch,
		Code:           "EvaluatorUnterminatedTimestampFormatPatternToken",
		Description:    "Time stamp format pattern contains unterminated token in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorUnterminatedTimestampFormatPatternToken: {
		ErrCode:        ErrEvaluatorUnterminatedTimestampFormatPatternToken,
		Code:           "EvaluatorInvalidTimestampFormatPatternToken",
		Description:    "Time stamp format pattern contains an invalid token in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorInvalidTimestampFormatPatternToken: {
		ErrCode:        ErrEvaluatorInvalidTimestampFormatPatternToken,
		Code:           "EvaluatorInvalidTimestampFormatPatternToken",
		Description:    "Time stamp format pattern contains an invalid token in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorInvalidTimestampFormatPatternSymbol: {
		ErrCode:        ErrEvaluatorInvalidTimestampFormatPatternSymbol,
		Code:           "EvaluatorInvalidTimestampFormatPatternSymbol",
		Description:    "Time stamp format pattern contains an invalid symbol in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorBindingDoesNotExist: {
		ErrCode:        ErrEvaluatorBindingDoesNotExist,
		Code:           "ErrEvaluatorBindingDoesNotExist",
		Description:    "A column name or a path provided does not exist in the SQL expression",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingHeaders: {
		ErrCode:        ErrMissingHeaders,
		Code:           "MissingHeaders",
		Description:    "Some headers in the query are missing from the file. Check the file and try again.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidColumnIndex: {
		ErrCode:        ErrInvalidColumnIndex,
		Code:           "InvalidColumnIndex",
		Description:    "The column index is invalid. Please check the service documentation and try again.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidDecompressedSize: {
		ErrCode:        ErrInvalidDecompressedSize,
		Code:           "XMinioInvalidDecompressedSize",
		Description:    "The data provided is unfit for decompression",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAddUserInvalidArgument: {
		ErrCode:        ErrAddUserInvalidArgument,
		Code:           "XMinioInvalidIAMCredentials",
		Description:    "User is not allowed to be same as admin access key",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrAdminAccountNotEligible: {
		ErrCode:        ErrAdminAccountNotEligible,
		Code:           "XMinioInvalidIAMCredentials",
		Description:    "The administrator key is not eligible for this operation",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrServiceAccountNotFound: {
		ErrCode:        ErrServiceAccountNotFound,
		Code:           "XMinioInvalidIAMCredentials",
		Description:    "The specified service account is not found",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrPostPolicyConditionInvalidFormat: {
		ErrCode:        ErrPostPolicyConditionInvalidFormat,
		Code:           "PostPolicyInvalidKeyName",
		Description:    "Invalid according to Policy: Policy Condition failed",
		HTTPStatusCode: http.StatusForbidden,
	},
	// Add your error structure here.
}

// IsS3Error check if the provided error is a specific s3 error.
func IsS3Error(err error, code ErrorCode) bool {
	e, ok := err.(Error)
	return ok && e.ErrCode == code
}

func (e errorCodeMap) toAPIErrWithErr(errCode ErrorCode, err error) Error {
	apiErr, ok := e[errCode]
	if !ok {
		apiErr = e[ErrInternalError]
	}
	if err != nil {
		apiErr.Description = fmt.Sprintf("%s (%s)", apiErr.Description, err)
	}
	return apiErr
}

func (e errorCodeMap) toAPIErr(errCode ErrorCode) Error {
	return e.toAPIErrWithErr(errCode, nil)
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %d => %s", e.Code, e.HTTPStatusCode, e.Description)
}

// GetAPIError provides API Error for input API error code.
func GetAPIError(code ErrorCode) Error {
	if apiErr, ok := errorCodes[code]; ok {
		return apiErr
	}
	return errorCodes.toAPIErr(ErrInternalError)
}

// ObjectError - error that linked to specific object.
type ObjectError struct {
	Err     error
	Object  string
	Version string
}

func (e ObjectError) Error() string {
	return fmt.Sprintf("%s (%s:%s)", e.Err, e.Object, e.Version)
}

// ObjectVersion get "object:version" string.
func (e ObjectError) ObjectVersion() string {
	return e.Object + ":" + e.Version
}
