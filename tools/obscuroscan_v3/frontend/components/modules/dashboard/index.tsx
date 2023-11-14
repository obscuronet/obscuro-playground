import React from "react";
import { CalendarDateRangePicker } from "@/components/date-range-picker";
import { CardHeader, CardTitle, CardContent, Card } from "@/components/ui/card";
import {
  LayersIcon,
  FileTextIcon,
  ReaderIcon,
  CubeIcon,
  RocketIcon,
} from "@radix-ui/react-icons";

import { RecentBatches } from "./recent-batches";
import { RecentTransactions } from "./recent-transactions";
import { Button } from "@/components/ui/button";
import { useTransactions } from "@/src/hooks/useTransactions";
import { useBatches } from "@/src/hooks/useBatches";
import TruncatedAddress from "../common/truncated-address";
import { useContracts } from "@/src/hooks/useContracts";
import { Skeleton } from "@/components/ui/skeleton";
import { RecentBlocks } from "./recent-blocks";
import { useBlocks } from "@/src/hooks/useBlocks";
import AnalyticsCard from "./analytics-card";
import Link from "next/link";

export default function Dashboard() {
  const { price, transactions, transactionCount } = useTransactions();
  const { contractCount } = useContracts();
  const { batches, latestBatch } = useBatches();
  const { blocks } = useBlocks();

  const DASHBOARD_DATA = [
    {
      title: "Ether Price",
      value: price?.ethereum?.usd,
      change: "+20.1% from last month",
      icon: RocketIcon,
    },
    {
      title: "Latest Batch",
      value: batches?.result?.Total,
      change: "+20.1% from last month",
      icon: LayersIcon,
    },
    {
      title: "Latest Rollup",
      value: latestBatch?.item?.l1Proof ? (
        <TruncatedAddress
          address={latestBatch?.item?.l1Proof}
          prefixLength={6}
          suffixLength={4}
        />
      ) : (
        "N/A"
      ),
      change: "+20.1% from last month",
      icon: CubeIcon,
    },
    {
      title: "Transactions",
      value: transactionCount?.count,
      change: "+20.1% from last month",
      icon: ReaderIcon,
    },
    {
      title: "Contracts",
      value: contractCount?.count,
      change: "+20.1% from last month",
      icon: FileTextIcon,
    },
  ];

  const RECENT_DATA = [
    {
      title: "Recent Blocks",
      data: blocks,
      component: <RecentBlocks blocks={blocks} />,
      goTo: "/blocks",
      className: "sm:col-span-1 md:col-span-6 lg:col-span-3",
    },
    {
      title: "Recent Batches",
      data: batches,
      component: <RecentBatches batches={batches} />,
      goTo: "/batches",
      className: "sm:col-span-1 md:col-span-3 lg:col-span-3",
    },
    {
      title: "Recent Transactions",
      data: transactions,
      component: <RecentTransactions transactions={transactions} />,
      goTo: "/transactions",
      className: "sm:col-span-1 md:col-span-3 lg:col-span-3",
    },
  ];

  return (
    <>
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Obscuroscan</h2>
      </div>
      <div className="grid gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
        {DASHBOARD_DATA.map((item: any, index) => (
          <AnalyticsCard key={index} item={item} />
        ))}
      </div>
      <div className="grid gap-4 md:grid-cols-6 lg:grid-cols-9">
        {RECENT_DATA.map((item: any, index) => (
          <Card key={index} className={item.className}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 p-3">
              <CardTitle>{item.title}</CardTitle>
              <Link
                href={{
                  pathname: item.goTo,
                }}
              >
                <Button variant="outline" size="sm">
                  View All
                </Button>
              </Link>
            </CardHeader>
            <CardContent className="p-3">
              {item.data ? (
                item.component
              ) : (
                <Skeleton className="w-full h-[200px] rounded-lg" />
              )}
            </CardContent>
          </Card>
        ))}
      </div>
    </>
  );
}