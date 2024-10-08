import HeadSeo from "@/src/components/head-seo";
import Layout from "@/src/components/layouts/default-layout";
import EmptyState from "@/src/components/modules/common/empty-state";
import { Button } from "@/src/components/ui/button";
import { siteMetadata } from "@/src/lib/siteMetadata";
import { useRouter } from "next/router";
import React from "react";

const AddressDetails = () => {
  const { push } = useRouter();

  return (
    <Layout>
      <EmptyState
        title="Address Details"
        description="Coming soon..."
        imageSrc="/assets/images/clock.png"
        imageAlt="Clock"
        action={<Button onClick={() => push("/")}>Go Home</Button>}
      />
    </Layout>
  );
};

export default AddressDetails;

export async function getServerSideProps(context: any) {
  return {
    props: {},
  };
}
