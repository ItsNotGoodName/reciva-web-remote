/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface HttpHTTPError {
  message: string;
}

export interface HttpPatchState {
  audio_source?: string;
  power?: boolean;
  preset?: number;
  volume?: number;
}

export interface ModelBuild {
  built_by: string;
  commit: string;
  date: string;
  release_url: string;
  summary: string;
  version: string;
}

export interface ModelPreset {
  /** TitleNew is the overridden title. */
  title_new: string;
  /** URL of the preset. */
  url: string;
  /** URLNew is the overridden URL. */
  url_new: string;
}

export interface ModelRadio {
  name: string;
  uuid: string;
}

export interface StatePreset {
  /** Number is the preset number. */
  number: number;
  /** Title of the preset. */
  title: string;
  /** TitleNew is the overridden title. */
  title_new: string;
  /** URL of the preset. */
  url: string;
  /** URLNew is the overridden URL. */
  url_new: string;
}

export interface StateState {
  /** AudioSource is the audio source. */
  audio_source: string;
  /** AudioSources is the list of available audio sources. */
  audio_sources: string[];
  /** IsMuted represents if the radio is muted. */
  is_muted: boolean;
  /** Metadata of the current playing stream. */
  metadata: string;
  /** ModelName is the model name of the device. */
  model_name: string;
  /** ModelNumber is the model number of the device. */
  model_number: string;
  /** Name of the radio. */
  name: string;
  /** Power represents if the radio is not in standby. */
  power: boolean;
  /** PresetNumber is the current preset that is playing. */
  preset_number: number;
  /** Presets of the radio. */
  presets: StatePreset[];
  /** Status is either playing, connecting, or stopped. */
  status: StateStatus;
  /** Title of the current playing stream. */
  title: string;
  /** TitleNew is the overridden title. */
  title_new: string;
  /** URL of the stream that is currently selected. */
  url: string;
  /** URLNew is the overridden URL. */
  url_new: string;
  /** UUID of the radio. */
  uuid: string;
  /** Volume of the radio. */
  volume: number;
}

export enum StateStatus {
  StatusConnecting = "Connecting",
  StatusPlaying = "Playing",
  StatusStopped = "Stopped",
  StatusUnknown = "",
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;

export interface FullRequestParams extends Omit<RequestInit, "body"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseFormat;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (securityData: SecurityDataType | null) => Promise<RequestParams | void> | RequestParams | void;
  customFetch?: typeof fetch;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown> extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "/api";
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private abortControllers = new Map<CancelToken, AbortController>();
  private customFetch = (...fetchParams: Parameters<typeof fetch>) => fetch(...fetchParams);

  private baseApiParams: RequestParams = {
    credentials: "same-origin",
    headers: {},
    redirect: "follow",
    referrerPolicy: "no-referrer",
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected encodeQueryParam(key: string, value: any) {
    const encodedKey = encodeURIComponent(key);
    return `${encodedKey}=${encodeURIComponent(typeof value === "number" ? value : `${value}`)}`;
  }

  protected addQueryParam(query: QueryParamsType, key: string) {
    return this.encodeQueryParam(key, query[key]);
  }

  protected addArrayQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];
    return value.map((v: any) => this.encodeQueryParam(key, v)).join("&");
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter((key) => "undefined" !== typeof query[key]);
    return keys
      .map((key) => (Array.isArray(query[key]) ? this.addArrayQueryParam(query, key) : this.addQueryParam(query, key)))
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string") ? JSON.stringify(input) : input,
    [ContentType.Text]: (input: any) => (input !== null && typeof input !== "string" ? JSON.stringify(input) : input),
    [ContentType.FormData]: (input: any) =>
      Object.keys(input || {}).reduce((formData, key) => {
        const property = input[key];
        formData.append(
          key,
          property instanceof Blob
            ? property
            : typeof property === "object" && property !== null
            ? JSON.stringify(property)
            : `${property}`,
        );
        return formData;
      }, new FormData()),
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  protected mergeRequestParams(params1: RequestParams, params2?: RequestParams): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected createAbortSignal = (cancelToken: CancelToken): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = async <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format,
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.baseApiParams.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];
    const responseFormat = format || requestParams.format;

    return this.customFetch(`${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`, {
      ...requestParams,
      headers: {
        ...(requestParams.headers || {}),
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
      },
      signal: cancelToken ? this.createAbortSignal(cancelToken) : requestParams.signal,
      body: typeof body === "undefined" || body === null ? null : payloadFormatter(body),
    }).then(async (response) => {
      const r = response as HttpResponse<T, E>;
      r.data = null as unknown as T;
      r.error = null as unknown as E;

      const data = !responseFormat
        ? r
        : await response[responseFormat]()
            .then((data) => {
              if (r.ok) {
                r.data = data;
              } else {
                r.error = data;
              }
              return r;
            })
            .catch((e) => {
              r.error = e;
              return r;
            });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
    });
  };
}

/**
 * @title Reciva Web Remote
 * @version 1.0
 * @baseUrl /api
 * @contact
 *
 * Control your legacy Reciva based internet radios (Crane, Grace Digital, Tangent, etc.) via web browser or REST API.
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  build = {
    /**
     * No description
     *
     * @tags build
     * @name BuildList
     * @summary Get build
     * @request GET:/build
     */
    buildList: (params: RequestParams = {}) =>
      this.request<ModelBuild, any>({
        path: `/build`,
        method: "GET",
        format: "json",
        ...params,
      }),
  };
  presets = {
    /**
     * No description
     *
     * @tags presets
     * @name PresetsList
     * @summary List presets
     * @request GET:/presets
     */
    presetsList: (params: RequestParams = {}) =>
      this.request<ModelPreset[], HttpHTTPError>({
        path: `/presets`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags presets
     * @name PresetsCreate
     * @summary Update preset
     * @request POST:/presets
     */
    presetsCreate: (preset: ModelPreset, params: RequestParams = {}) =>
      this.request<void, HttpHTTPError>({
        path: `/presets`,
        method: "POST",
        body: preset,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * No description
     *
     * @tags presets
     * @name PresetsDetail
     * @summary Get preset
     * @request GET:/presets/{url}
     */
    presetsDetail: (url: string, params: RequestParams = {}) =>
      this.request<ModelPreset[], HttpHTTPError>({
        path: `/presets/${url}`,
        method: "GET",
        format: "json",
        ...params,
      }),
  };
  radios = {
    /**
     * No description
     *
     * @tags radios
     * @name RadiosList
     * @summary List radios
     * @request GET:/radios
     */
    radiosList: (params: RequestParams = {}) =>
      this.request<ModelRadio[], any>({
        path: `/radios`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags radios
     * @name RadiosCreate
     * @summary Discover radios
     * @request POST:/radios
     */
    radiosCreate: (params: RequestParams = {}) =>
      this.request<void, HttpHTTPError>({
        path: `/radios`,
        method: "POST",
        ...params,
      }),

    /**
     * No description
     *
     * @tags radios
     * @name RadiosDetail
     * @summary Get radio
     * @request GET:/radios/{uuid}
     */
    radiosDetail: (uuid: string, params: RequestParams = {}) =>
      this.request<ModelRadio, HttpHTTPError>({
        path: `/radios/${uuid}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags radios
     * @name SubscriptionCreate
     * @summary Refresh radio subscription
     * @request POST:/radios/{uuid}/subscription
     */
    subscriptionCreate: (uuid: string, params: RequestParams = {}) =>
      this.request<void, HttpHTTPError>({
        path: `/radios/${uuid}/subscription`,
        method: "POST",
        ...params,
      }),

    /**
     * No description
     *
     * @tags radios
     * @name VolumeCreate
     * @summary Refresh radio volume
     * @request POST:/radios/{uuid}/volume
     */
    volumeCreate: (uuid: string, params: RequestParams = {}) =>
      this.request<void, HttpHTTPError>({
        path: `/radios/${uuid}/volume`,
        method: "POST",
        ...params,
      }),
  };
  states = {
    /**
     * No description
     *
     * @tags states
     * @name StatesList
     * @summary List states
     * @request GET:/states
     */
    statesList: (params: RequestParams = {}) =>
      this.request<StateState[], any>({
        path: `/states`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags states
     * @name StatesDetail
     * @summary Get state
     * @request GET:/states/{uuid}
     */
    statesDetail: (uuid: string, params: RequestParams = {}) =>
      this.request<StateState[], HttpHTTPError>({
        path: `/states/${uuid}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags states
     * @name StatesPartialUpdate
     * @summary Patch state
     * @request PATCH:/states/{uuid}
     */
    statesPartialUpdate: (uuid: string, state: HttpPatchState, params: RequestParams = {}) =>
      this.request<void, HttpHTTPError>({
        path: `/states/${uuid}`,
        method: "PATCH",
        body: state,
        type: ContentType.Json,
        ...params,
      }),
  };
}
