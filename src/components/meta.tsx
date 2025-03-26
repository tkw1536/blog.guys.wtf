import * as React from 'react';
import Head from "next/head";

interface MetaProps {
    title: string;
    url?: string;
    description?: string;
    image?: string;
    children?: React.ReactNode;
}

/** Meta generates meta-information to be included into the final page */
export default class Meta extends React.Component<MetaProps> {
    render() {
        const { title, url, description, image, children } = this.props;
        return <>
            <Head>
                <meta name="viewport" content="width=device-width, initial-scale=1" />

                {url && <>
                    <meta name="og:url" content={url} />
                    <meta property="twitter:url" content={url} />
                </>}

                <title>{title}</title>

                <meta itemProp="name" content={title} />
                <meta name="og:title" content={title} />
                <meta property="twitter:title" content={title} />

                {description && <>
                    <meta name="description" content={description} />
                    <meta itemProp="description" content={description} />
                    <meta name="og:description" content={description} />
                    <meta name="og:site_name" content={description} />
                    <meta property="twitter:description" content={description} />
                </>}

                {image && <>
                    <meta name="image" content={image} />
                    <meta itemProp="image" content={image} />
                    <meta name="og:image" content={image} />
                    <meta property="twitter:image" content={image} />

                    <meta property="twitter:card" content="summary_large_image" />
                </>}


                <meta name="og:type" content="website" />
            </Head>
            {children}
        </>;
    }
}