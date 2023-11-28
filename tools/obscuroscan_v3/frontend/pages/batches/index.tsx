import React from "react";
import { columns } from "@/src/components/modules/batches/columns";
import { DataTable } from "@/src/components/modules/common/data-table/data-table";
import Layout from "@/src/components/layouts/default-layout";
import { Metadata } from "next";
import { useBatchesService } from "@/src/services/useBatchesService";

export const metadata: Metadata = {
  title: "Batches",
  description: "A table of Batches.",
};

export default function Batches() {
  const { batches, refetchBatches, setNoPolling } = useBatchesService();
  const { BatchesData, Total } = batches?.result || {
    BatchesData: [],
    Total: 0,
  };

  React.useEffect(() => {
    setNoPolling(true);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <Layout>
      <div className="h-full flex-1 flex-col space-y-8 md:flex">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Batches</h2>
            <p className="text-sm text-muted-foreground">
              {Total} Batches found.
            </p>
          </div>
        </div>
        {BatchesData ? (
          <DataTable
            columns={columns}
            data={BatchesData}
            refetch={refetchBatches}
            total={+Total}
          />
        ) : (
          <p>Loading...</p>
        )}
      </div>
    </Layout>
  );
}
