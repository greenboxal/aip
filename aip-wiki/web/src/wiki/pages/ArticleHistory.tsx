import React from "react";

import {gql} from "../../__generated__";
import {useQuery} from "@apollo/client";

import ReactMarkdown, {Components} from 'react-markdown'
import remarkGfm from 'remark-gfm'
import remarkMath from "remark-math";
import rehypeSlug from "rehype-slug";
import rehypeKatex from "rehype-katex";
import rehypeSectionize from "@hbsnow/rehype-sectionize";
import rehypeAutolinkHeadings from "rehype-autolink-headings";

import {useLocation, useParams} from "react-router";
import {Box, Card, Icon, Link, Skeleton, Stack, Typography} from "@mui/material";
import LinkIcon from '@mui/icons-material/Link';
import Image from "mui-image";
import {useSearchParams} from "react-router-dom";

export const GET_PAGE = gql(/* GraphQL */ `
    query GetPage($id: String!) {
        Page(id: $id) {
            metadata {
                id
            }
            spec {
                title
                language
                voice
            }
            status {
                markdown
            }
        }
    }
`)


type ArticleViewProps = {
    slug: string,
    pageId: string,
}

const Article: React.FC<ArticleViewProps> = (props) => {
    const { loading, data } = useQuery(GET_PAGE, {
        variables: {
            id: props.pageId,
        }
    })

    if (loading) {
        return <article>
            <h1>{props.slug}</h1>
            <Skeleton variant="rectangular" />
        </article>
    }

    const components = {
        section: ({ node, children }) => (
            <Typography component="section" variant="body1">
                {children}
            </Typography>
        ),

        h1: ({ node, ...props }) => (<Typography variant="h1" {...props} />),
        h2: ({ node, ...props }) => (<Typography variant="h2" {...props} />),
        h3: ({ node, ...props }) => (<Typography variant="h3" {...props} />),
        h4: ({ node, ...props }) => (<Typography variant="h4" {...props} />),
        h5: ({ node, ...props }) => (<Typography variant="h5" {...props} />),
        h6: ({ node, ...props }) => (<Typography variant="h6" {...props} />),

        ArticleSectionHeading: ({ node }) => {
            return <LinkIcon />
        },

        a: ({ node, href, children }) => (
            <Link href={href}>
                {children}
            </Link>
        ),

        img: ({ node, src, alt, title }) => (
            <Card sx={{
                clear: 'right',
                float: 'right',
                maxWidth: '30%',
            }}>
                <Image fit="contain" src={src} alt={alt} title={title} />
            </Card>
        )
    } as Components

    return <Box component="article" sx={{ mt: 8 }}>
        <ReactMarkdown
            children={data.Page.status.markdown}
            components={components}
            remarkPlugins={[
                remarkGfm,
                remarkMath,
            ]}
            rehypePlugins={[
                rehypeKatex,
                rehypeSlug,
                rehypeSectionize,

                [rehypeAutolinkHeadings, {
                    behavior: 'append',
                    properties: {
                        className: 'anchor-link',
                    },
                    content: () => {
                        return {
                            type: 'element',
                            tagName: 'ArticleSectionHeading',
                            properties: {},
                            children: [],
                        }
                    },
                }],
            ]}
        />
    </Box>
}

const ArticlePage = () => {
    const params = useParams()
    const [search, _]  = useSearchParams();

    const slug = params["slug"]
    const pageId = search.get("pageId")

    return <Article slug={slug} pageId={pageId} />
}

export default ArticlePage