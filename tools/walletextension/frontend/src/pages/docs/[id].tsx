import Layout from "@/components/layouts/default-layout";
import { siteMetadata } from "@/lib/siteMetadata";
import { useRouter } from "next/router";
import React from "react";
import { ToastType } from "@repo/ui/lib/enums/toast";
import { showToast } from "@repo/ui/components/shared/use-toast";
import Spinner from "@repo/ui/components/shared/spinner";
import Custom404Error from "@repo/ui/components/common/404";

type Document = {
  title: string;
  subHeading: string;
  content: {
    heading: string;
    content: string[];
  }[];
};

const Document = () => {
  const { query } = useRouter();
  const { id } = query;

  const [document, setDocument] = React.useState<Document>({} as Document);
  const [loading, setLoading] = React.useState<boolean>(false);

  const getDocument = async () => {
    setLoading(true);
    try {
      const response = await fetch(`/docs/${id}.json`);
      const data = await response.json();
      const processedData = {
        title: data.title,
        subHeading: data.subHeading,
        content: data.content.map((item: any) => {
          return {
            heading: item.heading,
            content: item.content.map((paragraph: any) => {
              return paragraph.replace(
                /siteMetadata.email/g,
                siteMetadata.email
              );
            }),
          };
        }),
      };
      setDocument(processedData);
    } catch (error) {
      showToast(ToastType.DESTRUCTIVE, "Error fetching document");
    } finally {
      setLoading(false);
    }
  };

  React.useEffect(() => {
    if (id) {
      getDocument();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  return (
    <Layout>
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        {!loading ? (
          !document.title ? (
            <Custom404Error customPageTitle="Document" />
          ) : (
            <>
              <div className="mb-8 text-center">
                <h1 className="text-4xl font-extrabold mb-6">
                  {document.title}
                </h1>
                <p className="text-sm text-muted-foreground">
                  {document.subHeading}
                </p>
              </div>
              <div className="prose prose-lg prose-primary">
                {document.content &&
                  document.content.map((section, index) => (
                    <div key={index} className="mb-8">
                      <h2 className="mb-2">{section.heading}</h2>
                      {section.content &&
                        section.content.map((paragraph, index) => (
                          <div
                            key={index}
                            dangerouslySetInnerHTML={{ __html: paragraph }}
                          ></div>
                        ))}
                    </div>
                  ))}
              </div>
            </>
          )
        ) : (
          <Spinner />
        )}
      </div>
    </Layout>
  );
};

export default Document;
