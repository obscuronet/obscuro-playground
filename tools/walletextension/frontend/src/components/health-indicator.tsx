import React, { useEffect, useState } from "react";
import useGatewayService from "@/services/useGatewayService";
import { Badge, badgeVariants } from "@repo/ui/components/shared/badge";
import { Skeleton } from "@repo/ui/components/shared/skeleton";

const HealthIndicator = () => {
  const { getTestnetStatus } = useGatewayService();
  const [loading, setLoading] = useState(false);
  const [status, setStatus] = useState<boolean>();

  const testnetStatus = async () => {
    setLoading(true);
    try {
      const status = await getTestnetStatus();
      return status;
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let isMounted = true;

    // if overall health is true, status is set to true
    // if overall health is false but if the error includes [p2p], status is set to true

    testnetStatus().then((res) => {
      if (isMounted) {
        if (res?.result?.OverallHealth) {
          setStatus(true);
        } else if (
          res?.result?.Errors?.some((e: string) => e.includes("[p2p]"))
        ) {
          setStatus(true);
        } else {
          setStatus(false);
        }
      }
    });

    return () => {
      isMounted = false;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className="flex items-center space-x-2 border rounded-lg p-2">
      <h3 className="text-sm">Testnet Status: </h3>
      {loading ? (
        <Skeleton className="h-4 w-10" />
      ) : (
        <Badge
          variant={
            (status ? "success" : "destructive") as keyof typeof badgeVariants
          }
        >
          {status ? "Live" : "Down"}
        </Badge>
      )}
    </div>
  );
};

export default HealthIndicator;
