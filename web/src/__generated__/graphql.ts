/* eslint-disable */
import {TypedDocumentNode as DocumentNode} from '@graphql-typed-document-node/core';

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
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseAgentIdAgent>;
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
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseChannelIdChannel>;
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

export type Endpoint = {
    __typename?: 'Endpoint';
    id?: Maybe<Scalars['String']>;
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseEndpointIdEndpoint>;
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
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbContentAddressedResourceBaseImageIdImage>;
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
    metadata?: InputMaybe<IngithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseMemoryIdMemory>;
    parent_memory_id?: InputMaybe<Scalars['String']>;
    root_memory_id?: InputMaybe<Scalars['String']>;
};

export type InMemoryDatum = {
    cid?: InputMaybe<Scalars['String']>;
    data?: InputMaybe<Scalars['String']>;
    metadata?: InputMaybe<IngithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseMemoryDataIdMemoryData>;
};

export type IngithubcomgreenboxalaipaipcontrollerpkgfordforddbContentAddressedResourceBaseImageIdImage = {
    created_at?: InputMaybe<Scalars['DateTime']>;
    id?: InputMaybe<Scalars['String']>;
    kind?: InputMaybe<Scalars['String']>;
    namespace?: InputMaybe<Scalars['String']>;
    scope?: InputMaybe<Scalars['String']>;
    updated_at?: InputMaybe<Scalars['DateTime']>;
    version?: InputMaybe<Scalars['Int']>;
};

export type IngithubcomgreenboxalaipaipcontrollerpkgfordforddbContentAddressedResourceBasePageIdPage = {
    created_at?: InputMaybe<Scalars['DateTime']>;
    id?: InputMaybe<Scalars['String']>;
    kind?: InputMaybe<Scalars['String']>;
    namespace?: InputMaybe<Scalars['String']>;
    scope?: InputMaybe<Scalars['String']>;
    updated_at?: InputMaybe<Scalars['DateTime']>;
    version?: InputMaybe<Scalars['Int']>;
};

export type IngithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseMemoryDataIdMemoryData = {
    created_at?: InputMaybe<Scalars['DateTime']>;
    id?: InputMaybe<Scalars['String']>;
    kind?: InputMaybe<Scalars['String']>;
    namespace?: InputMaybe<Scalars['String']>;
    scope?: InputMaybe<Scalars['String']>;
    updated_at?: InputMaybe<Scalars['DateTime']>;
    version?: InputMaybe<Scalars['Int']>;
};

export type IngithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseMemoryIdMemory = {
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
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseJobIdJob>;
    spec?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgjobsJobSpec>;
    status?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgjobsJobStatus>;
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

export type Memory = {
    __typename?: 'Memory';
    branch_memory_id?: Maybe<Scalars['String']>;
    clock?: Maybe<Scalars['Int']>;
    data?: Maybe<MemoryDatum>;
    height?: Maybe<Scalars['Int']>;
    id?: Maybe<Scalars['String']>;
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseMemoryIdMemory>;
    parent_memory_id?: Maybe<Scalars['String']>;
    root_memory_id?: Maybe<Scalars['String']>;
};

export type MemoryDatum = {
    __typename?: 'MemoryDatum';
    cid?: Maybe<Scalars['String']>;
    data?: Maybe<Scalars['String']>;
    id?: Maybe<Scalars['String']>;
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseMemoryDataIdMemoryData>;
};

export type MemoryDatumFilter = {
    id?: InputMaybe<Scalars['String']>;
    ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
    q?: InputMaybe<Scalars['String']>;
};

export type MemoryDatumListMetadata = {
    __typename?: 'MemoryDatumListMetadata';
    count?: Maybe<Scalars['Int']>;
};

export type MemoryFilter = {
    id?: InputMaybe<Scalars['String']>;
    ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
    q?: InputMaybe<Scalars['String']>;
};

export type MemoryListMetadata = {
    __typename?: 'MemoryListMetadata';
    count?: Maybe<Scalars['Int']>;
};

export type MemorySegment = {
    __typename?: 'MemorySegment';
    id?: Maybe<Scalars['String']>;
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseMemorySegmentIdMemorySegment>;
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
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseMessageIdMessage>;
    reply_to_id?: Maybe<Scalars['String']>;
    text?: Maybe<Scalars['String']>;
    thread_id?: Maybe<Scalars['String']>;
};

export type MessageFilter = {
    id?: InputMaybe<Scalars['String']>;
    ids?: InputMaybe<Array<InputMaybe<Scalars['String']>>>;
    q?: InputMaybe<Scalars['String']>;
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
    supervisorListChildren?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgapissupervisorListChildrenResponse>;
    supervisorStartChild?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgapissupervisorStartChildResponse>;
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
    StringResourceID?: InputMaybe<Scalars['String']>;
};


export type MutationsWikiContentCachePutImageArgs = {
    metadata?: InputMaybe<IngithubcomgreenboxalaipaipcontrollerpkgfordforddbContentAddressedResourceBaseImageIdImage>;
    spec?: InputMaybe<IngithubcomgreenboxalaipaipwikipkgwikimodelsImageSpec>;
    status?: InputMaybe<IngithubcomgreenboxalaipaipwikipkgwikimodelsImageStatus>;
};


export type MutationsWikiContentCachePutPageArgs = {
    metadata?: InputMaybe<IngithubcomgreenboxalaipaipcontrollerpkgfordforddbContentAddressedResourceBasePageIdPage>;
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
    StringResourceID?: InputMaybe<Scalars['String']>;
};

export type Page = {
    __typename?: 'Page';
    id?: Maybe<Scalars['String']>;
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbContentAddressedResourceBasePageIdPage>;
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
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBasePipelineIdPipeline>;
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
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBasePortIdPort>;
    spec?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectivePortSpec>;
    status?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgcollectivePortStatus>;
};

export type PortBinding = {
    __typename?: 'PortBinding';
    id?: Maybe<Scalars['String']>;
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBasePortBindingIdPortBinding>;
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
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseProfileIdProfile>;
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
    Endpoint?: Maybe<Endpoint>;
    Image?: Maybe<Image>;
    Job?: Maybe<Job>;
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
    Stage?: Maybe<Stage>;
    Task?: Maybe<Task>;
    Team?: Maybe<Team>;
    _allAgentsMeta?: Maybe<AgentListMetadata>;
    _allChannelsMeta?: Maybe<ChannelListMetadata>;
    _allEndpointsMeta?: Maybe<EndpointListMetadata>;
    _allImagesMeta?: Maybe<ImageListMetadata>;
    _allJobsMeta?: Maybe<JobListMetadata>;
    _allMemoriesMeta?: Maybe<MemoryListMetadata>;
    _allMemoryDataMeta?: Maybe<MemoryDatumListMetadata>;
    _allMemorySegmentsMeta?: Maybe<MemorySegmentListMetadata>;
    _allMessagesMeta?: Maybe<MessageListMetadata>;
    _allPagesMeta?: Maybe<PageListMetadata>;
    _allPipelinesMeta?: Maybe<PipelineListMetadata>;
    _allPortBindingsMeta?: Maybe<PortBindingListMetadata>;
    _allPortsMeta?: Maybe<PortListMetadata>;
    _allProfilesMeta?: Maybe<ProfileListMetadata>;
    _allStagesMeta?: Maybe<StageListMetadata>;
    _allTasksMeta?: Maybe<TaskListMetadata>;
    _allTeamsMeta?: Maybe<TeamListMetadata>;
    allAgents?: Maybe<Array<Maybe<Agent>>>;
    allChannels?: Maybe<Array<Maybe<Channel>>>;
    allEndpoints?: Maybe<Array<Maybe<Endpoint>>>;
    allImages?: Maybe<Array<Maybe<Image>>>;
    allJobs?: Maybe<Array<Maybe<Job>>>;
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
    allStages?: Maybe<Array<Maybe<Stage>>>;
    allTasks?: Maybe<Array<Maybe<Task>>>;
    allTeams?: Maybe<Array<Maybe<Team>>>;
};


export type RootQueryAgentArgs = {
    id?: InputMaybe<Scalars['String']>;
};


export type RootQueryChannelArgs = {
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


export type RootQueryStageArgs = {
    id?: InputMaybe<Scalars['String']>;
};


export type RootQueryTaskArgs = {
    id?: InputMaybe<Scalars['String']>;
};


export type RootQueryTeamArgs = {
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

export type Stage = {
    __typename?: 'Stage';
    assigned_team?: Maybe<Scalars['String']>;
    depends_on?: Maybe<Array<Maybe<Scalars['String']>>>;
    id?: Maybe<Scalars['String']>;
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseStageIdStage>;
};

export type StageFilter = {
    id?: InputMaybe<Scalars['String']>;
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
    resourceChanged?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgapisgraphqlResourceEvent>;
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
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseTaskIdTask>;
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
    metadata?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseTeamIdTeam>;
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

export type GithubcomgreenboxalaipaipcontrollerpkgapisgraphqlResourceEvent = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgapisgraphqlResourceEvent';
    id?: Maybe<Scalars['String']>;
    payload?: Maybe<GithubcomgreenboxalaipaipcontrollerpkgapisgraphqlResourceEventPayload>;
    type?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgapisgraphqlResourceEventPayload = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgapisgraphqlResourceEventPayload';
    id?: Maybe<Scalars['String']>;
    ids?: Maybe<Array<Maybe<Scalars['String']>>>;
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

export type GithubcomgreenboxalaipaipcontrollerpkgapissupervisorListChildrenResponse = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgapissupervisorListChildrenResponse';
    children?: Maybe<Array<Maybe<Scalars['String']>>>;
    id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgapissupervisorStartChildResponse = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgapissupervisorStartChildResponse';
    id?: Maybe<Scalars['String']>;
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

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbContentAddressedResourceBaseImageIdImage = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbContentAddressedResourceBaseImageIDImage';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbContentAddressedResourceBasePageIdPage = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbContentAddressedResourceBasePageIDPage';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseAgentIdAgent = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseAgentIDAgent';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseChannelIdChannel = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseChannelIDChannel';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseEndpointIdEndpoint = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseEndpointIDEndpoint';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseJobIdJob = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseJobIDJob';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseMemoryDataIdMemoryData = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseMemoryDataIDMemoryData';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseMemoryIdMemory = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseMemoryIDMemory';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseMemorySegmentIdMemorySegment = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseMemorySegmentIDMemorySegment';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseMessageIdMessage = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseMessageIDMessage';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBasePipelineIdPipeline = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBasePipelineIDPipeline';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBasePortBindingIdPortBinding = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBasePortBindingIDPortBinding';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBasePortIdPort = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBasePortIDPort';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseProfileIdProfile = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseProfileIDProfile';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseStageIdStage = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseStageIDStage';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseTaskIdTask = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseTaskIDTask';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseTeamIdTeam = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordforddbResourceBaseTeamIDTeam';
    created_at?: Maybe<Scalars['DateTime']>;
    id?: Maybe<Scalars['String']>;
    kind?: Maybe<Scalars['String']>;
    namespace?: Maybe<Scalars['String']>;
    scope?: Maybe<Scalars['String']>;
    updated_at?: Maybe<Scalars['DateTime']>;
    version?: Maybe<Scalars['Int']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgfordreconciliationreconcilerInformation = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgfordreconciliationreconcilerInformation';
    id?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgjobsJobSpec = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgjobsJobSpec';
    handler?: Maybe<Scalars['String']>;
    id?: Maybe<Scalars['String']>;
    payload?: Maybe<Scalars['String']>;
};

export type GithubcomgreenboxalaipaipcontrollerpkgjobsJobStatus = {
    __typename?: 'githubcomgreenboxalaipaipcontrollerpkgjobsJobStatus';
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
    resourceKind: Scalars['String'];
}>;


export type SubSubscription = {
    __typename?: 'Subscriptions',
    resourceChanged?: {
        __typename?: 'githubcomgreenboxalaipaipcontrollerpkgapisgraphqlResourceEvent',
        type?: string | null,
        payload?: { __typename?: 'githubcomgreenboxalaipaipcontrollerpkgapisgraphqlResourceEventPayload', ids?: Array<string | null> | null } | null
    } | null
};


export const SubDocument = {
    "kind": "Document",
    "definitions": [{
        "kind": "OperationDefinition",
        "operation": "subscription",
        "name": {"kind": "Name", "value": "Sub"},
        "variableDefinitions": [{
            "kind": "VariableDefinition",
            "variable": {"kind": "Variable", "name": {"kind": "Name", "value": "resourceKind"}},
            "type": {"kind": "NonNullType", "type": {"kind": "NamedType", "name": {"kind": "Name", "value": "String"}}}
        }],
        "selectionSet": {
            "kind": "SelectionSet",
            "selections": [{
                "kind": "Field",
                "name": {"kind": "Name", "value": "resourceChanged"},
                "arguments": [{
                    "kind": "Argument",
                    "name": {"kind": "Name", "value": "resourceType"},
                    "value": {"kind": "Variable", "name": {"kind": "Name", "value": "resourceKind"}}
                }],
                "selectionSet": {
                    "kind": "SelectionSet",
                    "selections": [{"kind": "Field", "name": {"kind": "Name", "value": "type"}}, {
                        "kind": "Field",
                        "name": {"kind": "Name", "value": "payload"},
                        "selectionSet": {"kind": "SelectionSet", "selections": [{"kind": "Field", "name": {"kind": "Name", "value": "ids"}}]}
                    }]
                }
            }]
        }
    }]
} as unknown as DocumentNode<SubSubscription, SubSubscriptionVariables>;