import { cn } from "../../lib/utils";

function Skeleton({
  className,
  width,
  height,
  ...props
}: {
  className?: string;
  width?: number | string;
  height?: number | string;
}) {
  return (
    <span
      className={cn(
        "animate-pulse rounded-md bg-muted",
        width && `w-${width}`,
        height && `h-${height}`,
        className
      )}
      {...props}
    />
  );
}

export { Skeleton };
