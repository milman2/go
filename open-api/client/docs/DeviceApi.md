# DeviceApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**getDevice**](DeviceApi.md#getdevice) | **GET** /devices/{deviceId} | Get Device |



## getDevice

> Device getDevice(deviceId)

Get Device

### Example

```ts
import {
  Configuration,
  DeviceApi,
} from '';
import type { GetDeviceRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new DeviceApi();

  const body = {
    // string
    deviceId: deviceId_example,
  } satisfies GetDeviceRequest;

  try {
    const data = await api.getDevice(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **deviceId** | `string` |  | [Defaults to `undefined`] |

### Return type

[**Device**](Device.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | success |  -  |
| **404** | not-found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

