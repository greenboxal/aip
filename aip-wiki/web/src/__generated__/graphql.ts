/* eslint-disable */
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  /** The `DateTime` scalar type represents a DateTime. The DateTime is serialized as an RFC 3339 quoted string */
  DateTime: any;
};

export type Agent = {
  __typename?: 'Agent';
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseAgentIdAgent>;
  spec?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectiveAgentSpec>;
  status?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectiveAgentStatus>;
};

export type AgentFilter = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type AgentListMetadata = {
  __typename?: 'AgentListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type Channel = {
  __typename?: 'Channel';
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseChannelIdChannel>;
  subscribers?: Maybe<Array<Maybe<Scalars['String']>>>;
};

export type ChannelFilter = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type ChannelListMetadata = {
  __typename?: 'ChannelListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type Domain = {
  __typename?: 'Domain';
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseDomainIdDomain>;
  spec?: Maybe<GithubcomgreenboxalaipaipwikipkgwikimodelsDomainSpec>;
};

export type DomainFilter = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type DomainListMetadata = {
  __typename?: 'DomainListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type Endpoint = {
  __typename?: 'Endpoint';
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseEndpointIdEndpoint>;
  subscriptions?: Maybe<Array<Maybe<Scalars['String']>>>;
};

export type EndpointFilter = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type EndpointListMetadata = {
  __typename?: 'EndpointListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type Image = {
  __typename?: 'Image';
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbContentAddressedResourceBaseImageIdImage>;
  spec?: Maybe<GithubcomgreenboxalaipaipwikipkgwikimodelsImageSpec>;
  status?: Maybe<GithubcomgreenboxalaipaipwikipkgwikimodelsImageStatus>;
};

export type ImageFilter = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type ImageListMetadata = {
  __typename?: 'ImageListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type InMemory = {
  branch_memory_id?: InputMaybe<Scalars['String']>;
  clock?: InputMaybe<Scalars['Int']>;
  data?: InputMaybe<InMemoryDatum>;
  height?: InputMaybe<Scalars['Int']>;
  metadata?: InputMaybe<IngithubcomgreenboxalaipaipforddbpkgforddbResourceBaseMemoryIdMemory>;
  parent_memory_id?: InputMaybe<Scalars['String']>;
  root_memory_id?: InputMaybe<Scalars['String']>;
};

export type InMemoryDatum = {
  cid?: InputMaybe<Scalars['String']>;
  data?: InputMaybe<Scalars['String']>;
  metadata?: InputMaybe<IngithubcomgreenboxalaipaipforddbpkgforddbResourceBaseMemoryDataIdMemoryData>;
};

export type IngithubcomgreenboxalaipaipforddbpkgforddbContentAddressedResourceBaseImageIdImage = {
  created_at?: InputMaybe<Scalars['DateTime']>;
  id?: InputMaybe<Scalars['String']>;
  kind?: InputMaybe<Scalars['String']>;
  namespace?: InputMaybe<Scalars['String']>;
  scope?: InputMaybe<Scalars['String']>;
  updated_at?: InputMaybe<Scalars['DateTime']>;
  version?: InputMaybe<Scalars['Int']>;
};

export type IngithubcomgreenboxalaipaipforddbpkgforddbContentAddressedResourceBasePageIdPage = {
  created_at?: InputMaybe<Scalars['DateTime']>;
  id?: InputMaybe<Scalars['String']>;
  kind?: InputMaybe<Scalars['String']>;
  namespace?: InputMaybe<Scalars['String']>;
  scope?: InputMaybe<Scalars['String']>;
  updated_at?: InputMaybe<Scalars['DateTime']>;
  version?: InputMaybe<Scalars['Int']>;
};

export type IngithubcomgreenboxalaipaipforddbpkgforddbResourceBaseMemoryDataIdMemoryData = {
  created_at?: InputMaybe<Scalars['DateTime']>;
  id?: InputMaybe<Scalars['String']>;
  kind?: InputMaybe<Scalars['String']>;
  namespace?: InputMaybe<Scalars['String']>;
  scope?: InputMaybe<Scalars['String']>;
  updated_at?: InputMaybe<Scalars['DateTime']>;
  version?: InputMaybe<Scalars['Int']>;
};

export type IngithubcomgreenboxalaipaipforddbpkgforddbResourceBaseMemoryIdMemory = {
  created_at?: InputMaybe<Scalars['DateTime']>;
  id?: InputMaybe<Scalars['String']>;
  kind?: InputMaybe<Scalars['String']>;
  namespace?: InputMaybe<Scalars['String']>;
  scope?: InputMaybe<Scalars['String']>;
  updated_at?: InputMaybe<Scalars['DateTime']>;
  version?: InputMaybe<Scalars['Int']>;
};

export type IngithubcomgreenboxalaipaipwikipkgwikimodelsImageSpec = {
  path?: InputMaybe<Scalars['String']>;
  prompt?: InputMaybe<Scalars['String']>;
};

export type IngithubcomgreenboxalaipaipwikipkgwikimodelsImageStatus = {
  prompt?: InputMaybe<Scalars['String']>;
  url?: InputMaybe<Scalars['String']>;
};

export type IngithubcomgreenboxalaipaipwikipkgwikimodelsPageImage = {
  source?: InputMaybe<Scalars['String']>;
  title?: InputMaybe<Scalars['String']>;
};

export type IngithubcomgreenboxalaipaipwikipkgwikimodelsPageLink = {
  title?: InputMaybe<Scalars['String']>;
  to?: InputMaybe<Scalars['String']>;
};

export type IngithubcomgreenboxalaipaipwikipkgwikimodelsPageSpec = {
  base_page_id?: InputMaybe<Scalars['String']>;
  description?: InputMaybe<Scalars['String']>;
  format?: InputMaybe<Scalars['String']>;
  language?: InputMaybe<Scalars['String']>;
  layout?: InputMaybe<Scalars['String']>;
  title?: InputMaybe<Scalars['String']>;
  voice?: InputMaybe<Scalars['String']>;
};

export type IngithubcomgreenboxalaipaipwikipkgwikimodelsPageStatus = {
  html?: InputMaybe<Scalars['String']>;
  images?: InputMaybe<Array<InputMaybe<IngithubcomgreenboxalaipaipwikipkgwikimodelsPageImage>>>;
  links?: InputMaybe<Array<InputMaybe<IngithubcomgreenboxalaipaipwikipkgwikimodelsPageLink>>>;
  markdown?: InputMaybe<Scalars['String']>;
};

export type Job = {
  __typename?: 'Job';
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseJobIdJob>;
  spec?: Maybe<GithubcomgreenboxalaipaipsdkpkgjobsJobSpec>;
  status?: Maybe<GithubcomgreenboxalaipaipsdkpkgjobsJobStatus>;
};

export type JobFilter = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type JobListMetadata = {
  __typename?: 'JobListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type Layout = {
  __typename?: 'Layout';
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseLayoutIdLayout>;
  spec?: Maybe<GithubcomgreenboxalaipaipwikipkgwikimodelsLayoutSpec>;
};

export type LayoutFilter = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type LayoutListMetadata = {
  __typename?: 'LayoutListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type Memory = {
  __typename?: 'Memory';
  branch_memory_id?: Maybe<Scalars['String']>;
  clock?: Maybe<Scalars['Int']>;
  data?: Maybe<MemoryDatum>;
  height?: Maybe<Scalars['Int']>;
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseMemoryIdMemory>;
  parent_memory_id?: Maybe<Scalars['String']>;
  root_memory_id?: Maybe<Scalars['String']>;
};

export type MemoryDatum = {
  __typename?: 'MemoryDatum';
  cid?: Maybe<Scalars['String']>;
  data?: Maybe<Scalars['String']>;
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseMemoryDataIdMemoryData>;
};

export type MemoryDatumFilter = {
  data?: InputMaybe<Scalars['String']>;
  data_neq?: InputMaybe<Scalars['String']>;
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type MemoryDatumListMetadata = {
  __typename?: 'MemoryDatumListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type MemoryFilter = {
  branch_memory_id?: InputMaybe<Scalars['String']>;
  branch_memory_id_neq?: InputMaybe<Scalars['String']>;
  clock?: InputMaybe<Scalars['Int']>;
  clock_gt?: InputMaybe<Scalars['Int']>;
  clock_gte?: InputMaybe<Scalars['Int']>;
  clock_lt?: InputMaybe<Scalars['Int']>;
  clock_lte?: InputMaybe<Scalars['Int']>;
  clock_neq?: InputMaybe<Scalars['Int']>;
  height?: InputMaybe<Scalars['Int']>;
  height_gt?: InputMaybe<Scalars['Int']>;
  height_gte?: InputMaybe<Scalars['Int']>;
  height_lt?: InputMaybe<Scalars['Int']>;
  height_lte?: InputMaybe<Scalars['Int']>;
  height_neq?: InputMaybe<Scalars['Int']>;
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  parent_memory_id?: InputMaybe<Scalars['String']>;
  parent_memory_id_neq?: InputMaybe<Scalars['String']>;
  q?: InputMaybe<Scalars['String']>;
  root_memory_id?: InputMaybe<Scalars['String']>;
  root_memory_id_neq?: InputMaybe<Scalars['String']>;
};

export type MemoryListMetadata = {
  __typename?: 'MemoryListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type MemorySegment = {
  __typename?: 'MemorySegment';
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseMemorySegmentIdMemorySegment>;
};

export type MemorySegmentFilter = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type MemorySegmentListMetadata = {
  __typename?: 'MemorySegmentListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type Message = {
  __typename?: 'Message';
  channel?: Maybe<Scalars['String']>;
  from?: Maybe<Scalars['String']>;
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseMessageIdMessage>;
  reply_to_id?: Maybe<Scalars['String']>;
  role?: Maybe<Scalars['String']>;
  text?: Maybe<Scalars['String']>;
  thread_id?: Maybe<Scalars['String']>;
  username?: Maybe<Scalars['String']>;
};

export type MessageFilter = {
  channel?: InputMaybe<Scalars['String']>;
  channel_neq?: InputMaybe<Scalars['String']>;
  from?: InputMaybe<Scalars['String']>;
  from_neq?: InputMaybe<Scalars['String']>;
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
  reply_to_id?: InputMaybe<Scalars['String']>;
  reply_to_id_neq?: InputMaybe<Scalars['String']>;
  role?: InputMaybe<Scalars['String']>;
  role_neq?: InputMaybe<Scalars['String']>;
  text?: InputMaybe<Scalars['String']>;
  text_neq?: InputMaybe<Scalars['String']>;
  thread_id?: InputMaybe<Scalars['String']>;
  thread_id_neq?: InputMaybe<Scalars['String']>;
  username?: InputMaybe<Scalars['String']>;
  username_neq?: InputMaybe<Scalars['String']>;
};

export type MessageListMetadata = {
  __typename?: 'MessageListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type Mutations = {
  __typename?: 'Mutations';
  memlinkOneShotGetMemory?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgapismemorylinkOneShotGetMemoryResponse>;
  memlinkOneShotPutMemory?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgapismemorylinkOneShotPutMemoryResponse>;
  msnRouterJoinChannel?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectivemsnJoinChannelResponse>;
  msnRouterLeaveChannel?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectivemsnLeaveChannelResponse>;
  msnRouterPostMessage?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectivemsnPostMessageResponse>;
  supervisorListChildren?: Maybe<GithubcomgreenboxalaipaipsdkpkgapissupervisorListChildrenResponse>;
  supervisorStartChild?: Maybe<GithubcomgreenboxalaipaipsdkpkgapissupervisorStartChildResponse>;
  wikiContentCacheGetImage?: Maybe<Image>;
  wikiContentCacheGetPage?: Maybe<Page>;
  wikiContentCacheGetPageById?: Maybe<Page>;
  wikiContentCachePutImage?: Maybe<Image>;
  wikiContentCachePutPage?: Maybe<Page>;
  wikiEmpty?: Maybe<GithubcomgreenboxalaipaipwikipkgwikicmsEmptyRequest>;
  wikiImageGeneratorGetImage?: Maybe<GithubcomgreenboxalaipaipwikipkgwikimodelsImageStatus>;
  wikiPageGeneratorGeneratePage?: Maybe<Page>;
  wikiPageGeneratorGetPage?: Maybe<Array<Maybe<Scalars['Int']>>>;
  wikiPageManagerGenerateImage?: Maybe<Image>;
  wikiPageManagerGetImage?: Maybe<Image>;
  wikiPageManagerGetPage?: Maybe<Page>;
  wikiPageManagerGetPageById?: Maybe<Page>;
};


export type MutationsMemlinkOneShotGetMemoryArgs = {
  memory_id?: InputMaybe<Scalars['String']>;
};


export type MutationsMemlinkOneShotPutMemoryArgs = {
  new_memory?: InputMaybe<InMemoryDatum>;
  previous_memory?: InputMaybe<InMemory>;
};


export type MutationsMsnRouterJoinChannelArgs = {
  channel_id?: InputMaybe<Scalars['String']>;
  endpoint_id?: InputMaybe<Scalars['String']>;
};


export type MutationsMsnRouterLeaveChannelArgs = {
  channel_id?: InputMaybe<Scalars['String']>;
  endpoint_id?: InputMaybe<Scalars['String']>;
};


export type MutationsMsnRouterPostMessageArgs = {
  channel?: InputMaybe<Scalars['String']>;
  from?: InputMaybe<Scalars['String']>;
  text?: InputMaybe<Scalars['String']>;
};


export type MutationsSupervisorListChildrenArgs = {
  empty?: InputMaybe<Scalars['Boolean']>;
};


export type MutationsSupervisorStartChildArgs = {
  args?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  name?: InputMaybe<Scalars['String']>;
};


export type MutationsWikiContentCacheGetImageArgs = {
  path?: InputMaybe<Scalars['String']>;
  prompt?: InputMaybe<Scalars['String']>;
};


export type MutationsWikiContentCacheGetPageArgs = {
  base_page_id?: InputMaybe<Scalars['String']>;
  description?: InputMaybe<Scalars['String']>;
  format?: InputMaybe<Scalars['String']>;
  language?: InputMaybe<Scalars['String']>;
  layout?: InputMaybe<Scalars['String']>;
  title?: InputMaybe<Scalars['String']>;
  voice?: InputMaybe<Scalars['String']>;
};


export type MutationsWikiContentCacheGetPageByIdArgs = {
  arg?: InputMaybe<Scalars['String']>;
};


export type MutationsWikiContentCachePutImageArgs = {
  metadata?: InputMaybe<IngithubcomgreenboxalaipaipforddbpkgforddbContentAddressedResourceBaseImageIdImage>;
  spec?: InputMaybe<IngithubcomgreenboxalaipaipwikipkgwikimodelsImageSpec>;
  status?: InputMaybe<IngithubcomgreenboxalaipaipwikipkgwikimodelsImageStatus>;
};


export type MutationsWikiContentCachePutPageArgs = {
  metadata?: InputMaybe<IngithubcomgreenboxalaipaipforddbpkgforddbContentAddressedResourceBasePageIdPage>;
  spec?: InputMaybe<IngithubcomgreenboxalaipaipwikipkgwikimodelsPageSpec>;
  status?: InputMaybe<IngithubcomgreenboxalaipaipwikipkgwikimodelsPageStatus>;
};


export type MutationsWikiEmptyArgs = {
  empty?: InputMaybe<Scalars['String']>;
};


export type MutationsWikiImageGeneratorGetImageArgs = {
  path?: InputMaybe<Scalars['String']>;
  prompt?: InputMaybe<Scalars['String']>;
};


export type MutationsWikiPageGeneratorGeneratePageArgs = {
  base_page_id?: InputMaybe<Scalars['String']>;
  description?: InputMaybe<Scalars['String']>;
  format?: InputMaybe<Scalars['String']>;
  language?: InputMaybe<Scalars['String']>;
  layout?: InputMaybe<Scalars['String']>;
  title?: InputMaybe<Scalars['String']>;
  voice?: InputMaybe<Scalars['String']>;
};


export type MutationsWikiPageGeneratorGetPageArgs = {
  base_page_id?: InputMaybe<Scalars['String']>;
  description?: InputMaybe<Scalars['String']>;
  format?: InputMaybe<Scalars['String']>;
  language?: InputMaybe<Scalars['String']>;
  layout?: InputMaybe<Scalars['String']>;
  title?: InputMaybe<Scalars['String']>;
  voice?: InputMaybe<Scalars['String']>;
};


export type MutationsWikiPageManagerGenerateImageArgs = {
  path?: InputMaybe<Scalars['String']>;
  prompt?: InputMaybe<Scalars['String']>;
};


export type MutationsWikiPageManagerGetImageArgs = {
  path?: InputMaybe<Scalars['String']>;
  prompt?: InputMaybe<Scalars['String']>;
};


export type MutationsWikiPageManagerGetPageArgs = {
  base_page_id?: InputMaybe<Scalars['String']>;
  description?: InputMaybe<Scalars['String']>;
  format?: InputMaybe<Scalars['String']>;
  language?: InputMaybe<Scalars['String']>;
  layout?: InputMaybe<Scalars['String']>;
  title?: InputMaybe<Scalars['String']>;
  voice?: InputMaybe<Scalars['String']>;
};


export type MutationsWikiPageManagerGetPageByIdArgs = {
  arg?: InputMaybe<Scalars['String']>;
};

export type Page = {
  __typename?: 'Page';
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbContentAddressedResourceBasePageIdPage>;
  spec?: Maybe<GithubcomgreenboxalaipaipwikipkgwikimodelsPageSpec>;
  status?: Maybe<GithubcomgreenboxalaipaipwikipkgwikimodelsPageStatus>;
};

export type PageFilter = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type PageListMetadata = {
  __typename?: 'PageListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type Pipeline = {
  __typename?: 'Pipeline';
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBasePipelineIdPipeline>;
  spec?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectivePipelineSpec>;
};

export type PipelineFilter = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type PipelineListMetadata = {
  __typename?: 'PipelineListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type Port = {
  __typename?: 'Port';
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBasePortIdPort>;
  spec?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectivePortSpec>;
  status?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectivePortStatus>;
};

export type PortBinding = {
  __typename?: 'PortBinding';
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBasePortBindingIdPortBinding>;
  spec?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectivePortBindingSpec>;
  status?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectivePortBindingStatus>;
};

export type PortBindingFilter = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type PortBindingListMetadata = {
  __typename?: 'PortBindingListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type PortFilter = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type PortListMetadata = {
  __typename?: 'PortListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type Profile = {
  __typename?: 'Profile';
  aptitudes?: Maybe<Array<Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectiveProfileAptitude>>>;
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseProfileIdProfile>;
  spec?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectiveProfileSpec>;
};

export type ProfileFilter = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type ProfileListMetadata = {
  __typename?: 'ProfileListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type RootQuery = {
  __typename?: 'RootQuery';
  Agent?: Maybe<Agent>;
  Channel?: Maybe<Channel>;
  Domain?: Maybe<Domain>;
  Endpoint?: Maybe<Endpoint>;
  Image?: Maybe<Image>;
  Job?: Maybe<Job>;
  Layout?: Maybe<Layout>;
  Memory?: Maybe<Memory>;
  MemoryDatum?: Maybe<MemoryDatum>;
  MemorySegment?: Maybe<MemorySegment>;
  Message?: Maybe<Message>;
  Page?: Maybe<Page>;
  Pipeline?: Maybe<Pipeline>;
  Port?: Maybe<Port>;
  PortBinding?: Maybe<PortBinding>;
  Profile?: Maybe<Profile>;
  Reconciler?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordreconciliationreconcilerInformation>;
  Span?: Maybe<Span>;
  Stage?: Maybe<Stage>;
  Task?: Maybe<Task>;
  Team?: Maybe<Team>;
  Trace?: Maybe<Trace>;
  _allAgentsMeta?: Maybe<AgentListMetadata>;
  _allChannelsMeta?: Maybe<ChannelListMetadata>;
  _allDomainsMeta?: Maybe<DomainListMetadata>;
  _allEndpointsMeta?: Maybe<EndpointListMetadata>;
  _allImagesMeta?: Maybe<ImageListMetadata>;
  _allJobsMeta?: Maybe<JobListMetadata>;
  _allLayoutsMeta?: Maybe<LayoutListMetadata>;
  _allMemoriesMeta?: Maybe<MemoryListMetadata>;
  _allMemoryDataMeta?: Maybe<MemoryDatumListMetadata>;
  _allMemorySegmentsMeta?: Maybe<MemorySegmentListMetadata>;
  _allMessagesMeta?: Maybe<MessageListMetadata>;
  _allPagesMeta?: Maybe<PageListMetadata>;
  _allPipelinesMeta?: Maybe<PipelineListMetadata>;
  _allPortBindingsMeta?: Maybe<PortBindingListMetadata>;
  _allPortsMeta?: Maybe<PortListMetadata>;
  _allProfilesMeta?: Maybe<ProfileListMetadata>;
  _allSpansMeta?: Maybe<SpanListMetadata>;
  _allStagesMeta?: Maybe<StageListMetadata>;
  _allTasksMeta?: Maybe<TaskListMetadata>;
  _allTeamsMeta?: Maybe<TeamListMetadata>;
  _allTracesMeta?: Maybe<TraceListMetadata>;
  allAgents?: Maybe<Array<Maybe<Agent>>>;
  allChannels?: Maybe<Array<Maybe<Channel>>>;
  allDomains?: Maybe<Array<Maybe<Domain>>>;
  allEndpoints?: Maybe<Array<Maybe<Endpoint>>>;
  allImages?: Maybe<Array<Maybe<Image>>>;
  allJobs?: Maybe<Array<Maybe<Job>>>;
  allLayouts?: Maybe<Array<Maybe<Layout>>>;
  allMemories?: Maybe<Array<Maybe<Memory>>>;
  allMemoryData?: Maybe<Array<Maybe<MemoryDatum>>>;
  allMemorySegments?: Maybe<Array<Maybe<MemorySegment>>>;
  allMessages?: Maybe<Array<Maybe<Message>>>;
  allPages?: Maybe<Array<Maybe<Page>>>;
  allPipelines?: Maybe<Array<Maybe<Pipeline>>>;
  allPortBindings?: Maybe<Array<Maybe<PortBinding>>>;
  allPorts?: Maybe<Array<Maybe<Port>>>;
  allProfiles?: Maybe<Array<Maybe<Profile>>>;
  allReconcilers?: Maybe<Array<Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordreconciliationreconcilerInformation>>>;
  allSpans?: Maybe<Array<Maybe<Span>>>;
  allStages?: Maybe<Array<Maybe<Stage>>>;
  allTasks?: Maybe<Array<Maybe<Task>>>;
  allTeams?: Maybe<Array<Maybe<Team>>>;
  allTraces?: Maybe<Array<Maybe<Trace>>>;
};


export type RootQueryAgentArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryChannelArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryDomainArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryEndpointArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryImageArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryJobArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryLayoutArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryMemoryArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryMemoryDatumArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryMemorySegmentArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryMessageArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryPageArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryPipelineArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryPortArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryPortBindingArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryProfileArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryReconcilerArgs = {
  id: Scalars['String'];
};


export type RootQuerySpanArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryStageArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryTaskArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryTeamArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQueryTraceArgs = {
  id?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllAgentsMetaArgs = {
  filter?: InputMaybe<AgentFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllChannelsMetaArgs = {
  filter?: InputMaybe<ChannelFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllDomainsMetaArgs = {
  filter?: InputMaybe<DomainFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllEndpointsMetaArgs = {
  filter?: InputMaybe<EndpointFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllImagesMetaArgs = {
  filter?: InputMaybe<ImageFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllJobsMetaArgs = {
  filter?: InputMaybe<JobFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllLayoutsMetaArgs = {
  filter?: InputMaybe<LayoutFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllMemoriesMetaArgs = {
  filter?: InputMaybe<MemoryFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllMemoryDataMetaArgs = {
  filter?: InputMaybe<MemoryDatumFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllMemorySegmentsMetaArgs = {
  filter?: InputMaybe<MemorySegmentFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllMessagesMetaArgs = {
  filter?: InputMaybe<MessageFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllPagesMetaArgs = {
  filter?: InputMaybe<PageFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllPipelinesMetaArgs = {
  filter?: InputMaybe<PipelineFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllPortBindingsMetaArgs = {
  filter?: InputMaybe<PortBindingFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllPortsMetaArgs = {
  filter?: InputMaybe<PortFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllProfilesMetaArgs = {
  filter?: InputMaybe<ProfileFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllSpansMetaArgs = {
  filter?: InputMaybe<SpanFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllStagesMetaArgs = {
  filter?: InputMaybe<StageFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllTasksMetaArgs = {
  filter?: InputMaybe<TaskFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllTeamsMetaArgs = {
  filter?: InputMaybe<TeamFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQuery_AllTracesMetaArgs = {
  filter?: InputMaybe<TraceFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllAgentsArgs = {
  filter?: InputMaybe<AgentFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllChannelsArgs = {
  filter?: InputMaybe<ChannelFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllDomainsArgs = {
  filter?: InputMaybe<DomainFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllEndpointsArgs = {
  filter?: InputMaybe<EndpointFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllImagesArgs = {
  filter?: InputMaybe<ImageFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllJobsArgs = {
  filter?: InputMaybe<JobFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllLayoutsArgs = {
  filter?: InputMaybe<LayoutFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllMemoriesArgs = {
  filter?: InputMaybe<MemoryFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllMemoryDataArgs = {
  filter?: InputMaybe<MemoryDatumFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllMemorySegmentsArgs = {
  filter?: InputMaybe<MemorySegmentFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllMessagesArgs = {
  filter?: InputMaybe<MessageFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllPagesArgs = {
  filter?: InputMaybe<PageFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllPipelinesArgs = {
  filter?: InputMaybe<PipelineFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllPortBindingsArgs = {
  filter?: InputMaybe<PortBindingFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllPortsArgs = {
  filter?: InputMaybe<PortFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllProfilesArgs = {
  filter?: InputMaybe<ProfileFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllReconcilersArgs = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
};


export type RootQueryAllSpansArgs = {
  filter?: InputMaybe<SpanFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllStagesArgs = {
  filter?: InputMaybe<StageFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllTasksArgs = {
  filter?: InputMaybe<TaskFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllTeamsArgs = {
  filter?: InputMaybe<TeamFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};


export type RootQueryAllTracesArgs = {
  filter?: InputMaybe<TraceFilter>;
  page?: InputMaybe<Scalars['Int']>;
  perPage?: InputMaybe<Scalars['Int']>;
  sortField?: InputMaybe<Scalars['String']>;
  sortOrder?: InputMaybe<Scalars['String']>;
};

export type Span = {
  __typename?: 'Span';
  attributes?: Maybe<Array<Maybe<GithubcomgreenboxalaipaipforddbpkgtracingSpanAttribute>>>;
  completed_at?: Maybe<Scalars['DateTime']>;
  duration?: Maybe<Scalars['Int']>;
  id?: Maybe<Scalars['String']>;
  inner_span_ids?: Maybe<Array<Maybe<Scalars['String']>>>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseSpanIdSpan>;
  name?: Maybe<Scalars['String']>;
  parent_id?: Maybe<Scalars['String']>;
  started_at?: Maybe<Scalars['DateTime']>;
  trace_id?: Maybe<Scalars['String']>;
};

export type SpanFilter = {
  duration?: InputMaybe<Scalars['Int']>;
  duration_gt?: InputMaybe<Scalars['Int']>;
  duration_gte?: InputMaybe<Scalars['Int']>;
  duration_lt?: InputMaybe<Scalars['Int']>;
  duration_lte?: InputMaybe<Scalars['Int']>;
  duration_neq?: InputMaybe<Scalars['Int']>;
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  name?: InputMaybe<Scalars['String']>;
  name_neq?: InputMaybe<Scalars['String']>;
  parent_id?: InputMaybe<Scalars['String']>;
  parent_id_neq?: InputMaybe<Scalars['String']>;
  q?: InputMaybe<Scalars['String']>;
  trace_id?: InputMaybe<Scalars['String']>;
  trace_id_neq?: InputMaybe<Scalars['String']>;
};

export type SpanListMetadata = {
  __typename?: 'SpanListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type Stage = {
  __typename?: 'Stage';
  assigned_team?: Maybe<Scalars['String']>;
  depends_on?: Maybe<Array<Maybe<Scalars['String']>>>;
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseStageIdStage>;
};

export type StageFilter = {
  assigned_team?: InputMaybe<Scalars['String']>;
  assigned_team_neq?: InputMaybe<Scalars['String']>;
  id?: InputMaybe<Scalars['String']>;
  id_neq?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type StageListMetadata = {
  __typename?: 'StageListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type Subscriptions = {
  __typename?: 'Subscriptions';
  realTimeEvents?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectivemsnEvent>;
  resourceChanged?: Maybe<GithubcomgreenboxalaipaipsdkpkgapisgraphqlResourceEvent>;
};


export type SubscriptionsRealTimeEventsArgs = {
  endpoint: Scalars['String'];
};


export type SubscriptionsResourceChangedArgs = {
  resourceType: Scalars['String'];
};

export type Task = {
  __typename?: 'Task';
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseTaskIdTask>;
  spec?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectiveTaskSpec>;
  status?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectiveTaskStatus>;
};

export type TaskFilter = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type TaskListMetadata = {
  __typename?: 'TaskListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type Team = {
  __typename?: 'Team';
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseTeamIdTeam>;
  spec?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectiveTeamSpec>;
};

export type TeamFilter = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  q?: InputMaybe<Scalars['String']>;
};

export type TeamListMetadata = {
  __typename?: 'TeamListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type Trace = {
  __typename?: 'Trace';
  completed_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  metadata?: Maybe<GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseTraceIdTrace>;
  name?: Maybe<Scalars['String']>;
  root_span_id?: Maybe<Scalars['String']>;
  span_ids?: Maybe<Array<Maybe<Scalars['String']>>>;
  started_at?: Maybe<Scalars['DateTime']>;
};

export type TraceFilter = {
  id?: InputMaybe<Scalars['String']>;
  ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
  name?: InputMaybe<Scalars['String']>;
  name_neq?: InputMaybe<Scalars['String']>;
  q?: InputMaybe<Scalars['String']>;
  root_span_id?: InputMaybe<Scalars['String']>;
  root_span_id_neq?: InputMaybe<Scalars['String']>;
};

export type TraceListMetadata = {
  __typename?: 'TraceListMetadata';
  count?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgapismemorylinkOneShotGetMemoryResponse = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgapismemorylinkOneShotGetMemoryResponse';
  id?: Maybe<Scalars['String']>;
  memory?: Maybe<Memory>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgapismemorylinkOneShotPutMemoryResponse = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgapismemorylinkOneShotPutMemoryResponse';
  id?: Maybe<Scalars['String']>;
  new_memory?: Maybe<Memory>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectiveAgentSpec = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectiveAgentSpec';
  extra_args?: Maybe<Array<Maybe<Scalars['String']>>>;
  given_name?: Maybe<Scalars['String']>;
  id?: Maybe<Scalars['String']>;
  port_id?: Maybe<Scalars['String']>;
  profile_id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectiveAgentStatus = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectiveAgentStatus';
  id?: Maybe<Scalars['String']>;
  last_error?: Maybe<Scalars['String']>;
  state?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectivePipelineSpec = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectivePipelineSpec';
  id?: Maybe<Scalars['String']>;
  stages?: Maybe<Array<Maybe<Stage>>>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectivePortBindingSpec = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectivePortBindingSpec';
  agent_id?: Maybe<Scalars['String']>;
  id?: Maybe<Scalars['String']>;
  port_id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectivePortBindingStatus = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectivePortBindingStatus';
  id?: Maybe<Scalars['String']>;
  last_error?: Maybe<Scalars['String']>;
  state?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectivePortSpec = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectivePortSpec';
  empty?: Maybe<Scalars['Boolean']>;
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectivePortStatus = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectivePortStatus';
  id?: Maybe<Scalars['String']>;
  last_error?: Maybe<Scalars['String']>;
  state?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectiveProfileAptitude = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectiveProfileAptitude';
  description?: Maybe<Scalars['String']>;
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectiveProfileSpec = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectiveProfileSpec';
  directive?: Maybe<Scalars['String']>;
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectiveTaskPhaseStatus = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectiveTaskPhaseStatus';
  id?: Maybe<Scalars['String']>;
  stage_id?: Maybe<Scalars['String']>;
  state?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectiveTaskSpec = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectiveTaskSpec';
  description?: Maybe<Scalars['String']>;
  id?: Maybe<Scalars['String']>;
  output_stage_id?: Maybe<Scalars['String']>;
  pipeline_id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectiveTaskStatus = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectiveTaskStatus';
  id?: Maybe<Scalars['String']>;
  phase?: Maybe<Scalars['String']>;
  phases?: Maybe<Array<Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectiveTaskPhaseStatus>>>;
  state?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectiveTeamSpec = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectiveTeamSpec';
  id?: Maybe<Scalars['String']>;
  manager?: Maybe<Scalars['String']>;
  members?: Maybe<Array<Maybe<Scalars['String']>>>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectivemsnEvent = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectivemsnEvent';
  id?: Maybe<Scalars['String']>;
  message_event?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectivemsnMessageEvent>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectivemsnJoinChannelResponse = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectivemsnJoinChannelResponse';
  channel?: Maybe<Channel>;
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectivemsnLeaveChannelResponse = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectivemsnLeaveChannelResponse';
  channel?: Maybe<Channel>;
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectivemsnMessageEvent = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectivemsnMessageEvent';
  id?: Maybe<Scalars['String']>;
  message?: Maybe<Message>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgcollectivemsnPostMessageResponse = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgcollectivemsnPostMessageResponse';
  id?: Maybe<Scalars['String']>;
  message?: Maybe<Message>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordreconciliationreconcilerInformation = {
  __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordreconciliationreconcilerInformation';
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbContentAddressedResourceBaseImageIdImage = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbContentAddressedResourceBaseImageIDImage';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbContentAddressedResourceBasePageIdPage = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbContentAddressedResourceBasePageIDPage';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseAgentIdAgent = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBaseAgentIDAgent';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseChannelIdChannel = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBaseChannelIDChannel';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseDomainIdDomain = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBaseDomainIDDomain';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseEndpointIdEndpoint = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBaseEndpointIDEndpoint';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseJobIdJob = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBaseJobIDJob';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseLayoutIdLayout = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBaseLayoutIDLayout';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseMemoryDataIdMemoryData = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBaseMemoryDataIDMemoryData';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseMemoryIdMemory = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBaseMemoryIDMemory';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseMemorySegmentIdMemorySegment = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBaseMemorySegmentIDMemorySegment';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseMessageIdMessage = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBaseMessageIDMessage';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBasePipelineIdPipeline = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBasePipelineIDPipeline';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBasePortBindingIdPortBinding = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBasePortBindingIDPortBinding';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBasePortIdPort = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBasePortIDPort';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseProfileIdProfile = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBaseProfileIDProfile';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseSpanIdSpan = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBaseSpanIDSpan';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseStageIdStage = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBaseStageIDStage';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseTaskIdTask = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBaseTaskIDTask';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseTeamIdTeam = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBaseTeamIDTeam';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgforddbResourceBaseTraceIdTrace = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbResourceBaseTraceIDTrace';
  created_at?: Maybe<Scalars['DateTime']>;
  id?: Maybe<Scalars['String']>;
  kind?: Maybe<Scalars['String']>;
  namespace?: Maybe<Scalars['String']>;
  scope?: Maybe<Scalars['String']>;
  updated_at?: Maybe<Scalars['DateTime']>;
  version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipforddbpkgtracingSpanAttribute = {
  __typename?: 'githubcomgreenboxalaipaipforddbpkgtracingSpanAttribute';
  id?: Maybe<Scalars['String']>;
  key?: Maybe<Scalars['String']>;
  value?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipsdkpkgapisgraphqlResourceEvent = {
  __typename?: 'githubcomgreenboxalaipaipsdkpkgapisgraphqlResourceEvent';
  id?: Maybe<Scalars['String']>;
  payload?: Maybe<GithubcomgreenboxalaipaipsdkpkgapisgraphqlResourceEventPayload>;
  type?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipsdkpkgapisgraphqlResourceEventPayload = {
  __typename?: 'githubcomgreenboxalaipaipsdkpkgapisgraphqlResourceEventPayload';
  id?: Maybe<Scalars['String']>;
  ids?: Maybe<Array<Maybe<Scalars['String']>>>;
};

export type GithubcomgreenboxalaipaipsdkpkgapissupervisorListChildrenResponse = {
  __typename?: 'githubcomgreenboxalaipaipsdkpkgapissupervisorListChildrenResponse';
  children?: Maybe<Array<Maybe<Scalars['String']>>>;
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipsdkpkgapissupervisorStartChildResponse = {
  __typename?: 'githubcomgreenboxalaipaipsdkpkgapissupervisorStartChildResponse';
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipsdkpkgjobsJobSpec = {
  __typename?: 'githubcomgreenboxalaipaipsdkpkgjobsJobSpec';
  handler?: Maybe<Scalars['String']>;
  id?: Maybe<Scalars['String']>;
  payload?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipsdkpkgjobsJobStatus = {
  __typename?: 'githubcomgreenboxalaipaipsdkpkgjobsJobStatus';
  id?: Maybe<Scalars['String']>;
  last_error?: Maybe<Scalars['String']>;
  progress?: Maybe<Scalars['String']>;
  result?: Maybe<Scalars['String']>;
  state?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipwikipkgwikicmsEmptyRequest = {
  __typename?: 'githubcomgreenboxalaipaipwikipkgwikicmsEmptyRequest';
  empty?: Maybe<Scalars['String']>;
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipwikipkgwikimodelsDomainSpec = {
  __typename?: 'githubcomgreenboxalaipaipwikipkgwikimodelsDomainSpec';
  host?: Maybe<Scalars['String']>;
  id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipwikipkgwikimodelsImageSpec = {
  __typename?: 'githubcomgreenboxalaipaipwikipkgwikimodelsImageSpec';
  id?: Maybe<Scalars['String']>;
  path?: Maybe<Scalars['String']>;
  prompt?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipwikipkgwikimodelsImageStatus = {
  __typename?: 'githubcomgreenboxalaipaipwikipkgwikimodelsImageStatus';
  id?: Maybe<Scalars['String']>;
  prompt?: Maybe<Scalars['String']>;
  url?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipwikipkgwikimodelsLayoutSpec = {
  __typename?: 'githubcomgreenboxalaipaipwikipkgwikimodelsLayoutSpec';
  host?: Maybe<Scalars['String']>;
  id?: Maybe<Scalars['String']>;
  layout?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipwikipkgwikimodelsPageImage = {
  __typename?: 'githubcomgreenboxalaipaipwikipkgwikimodelsPageImage';
  id?: Maybe<Scalars['String']>;
  source?: Maybe<Scalars['String']>;
  title?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipwikipkgwikimodelsPageLink = {
  __typename?: 'githubcomgreenboxalaipaipwikipkgwikimodelsPageLink';
  id?: Maybe<Scalars['String']>;
  title?: Maybe<Scalars['String']>;
  to?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipwikipkgwikimodelsPageSpec = {
  __typename?: 'githubcomgreenboxalaipaipwikipkgwikimodelsPageSpec';
  base_page_id?: Maybe<Scalars['String']>;
  description?: Maybe<Scalars['String']>;
  format?: Maybe<Scalars['String']>;
  id?: Maybe<Scalars['String']>;
  language?: Maybe<Scalars['String']>;
  layout?: Maybe<Scalars['String']>;
  title?: Maybe<Scalars['String']>;
  voice?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipwikipkgwikimodelsPageStatus = {
  __typename?: 'githubcomgreenboxalaipaipwikipkgwikimodelsPageStatus';
  html?: Maybe<Scalars['String']>;
  id?: Maybe<Scalars['String']>;
  images?: Maybe<Array<Maybe<GithubcomgreenboxalaipaipwikipkgwikimodelsPageImage>>>;
  links?: Maybe<Array<Maybe<GithubcomgreenboxalaipaipwikipkgwikimodelsPageLink>>>;
  markdown?: Maybe<Scalars['String']>;
};

export type SubSubscriptionVariables = Exact<{
  resourceType: Scalars['String'];
}>;


export type SubSubscription = { __typename?: 'Subscriptions', resourceChanged?: { __typename?: 'githubcomgreenboxalaipaipsdkpkgapisgraphqlResourceEvent', type?: string | null, payload?: { __typename?: 'githubcomgreenboxalaipaipsdkpkgapisgraphqlResourceEventPayload', ids?: Array<string | null> | null } | null } | null };

export type GetPageQueryVariables = Exact<{
  id: Scalars['String'];
}>;


export type GetPageQuery = { __typename?: 'RootQuery', Page?: { __typename?: 'Page', metadata?: { __typename?: 'githubcomgreenboxalaipaipforddbpkgforddbContentAddressedResourceBasePageIDPage', id?: string | null } | null, spec?: { __typename?: 'githubcomgreenboxalaipaipwikipkgwikimodelsPageSpec', title?: string | null, language?: string | null, voice?: string | null } | null, status?: { __typename?: 'githubcomgreenboxalaipaipwikipkgwikimodelsPageStatus', markdown?: string | null } | null } | null };


export const SubDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"subscription","name":{"kind":"Name","value":"Sub"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"resourceType"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"resourceChanged"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"resourceType"},"value":{"kind":"Variable","name":{"kind":"Name","value":"resourceType"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"type"}},{"kind":"Field","name":{"kind":"Name","value":"payload"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"ids"}}]}}]}}]}}]} as unknown as DocumentNode<SubSubscription, SubSubscriptionVariables>;
export const GetPageDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetPage"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"Page"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"metadata"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"spec"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"language"}},{"kind":"Field","name":{"kind":"Name","value":"voice"}}]}},{"kind":"Field","name":{"kind":"Name","value":"status"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"markdown"}}]}}]}}]}}]} as unknown as DocumentNode<GetPageQuery, GetPageQueryVariables>;