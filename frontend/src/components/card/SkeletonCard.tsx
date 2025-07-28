import {Skeleton} from "@/components/ui/skeleton.tsx";
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from "@/components/ui/card.tsx";
import {Separator} from "@/components/ui/separator.tsx";

export default async function SkeletonCard({
  title
}: {
  title?: string
}) {
  return (
    <Card className="mb-2 overflow-hidden">
      <CardHeader className="pb-0">
        {title ?
          <CardTitle>{title}</CardTitle>
        :
          <Skeleton className="h-6 w-48"/>
        }
        <CardDescription/>
      </CardHeader>
      <CardContent className="p-0">
        <Separator className="mt-4"/>
        <div className="w-full relative">
          {[0,1,2].map(key =>
            <div key={key} className="containers tx-row-border">
              <div className="left">
                <Skeleton className="date h-6 w-32"/>
                <Skeleton className="remoteName h-6 w-40"/>
                <Skeleton className="purpose h-6 w-48"/>
              </div>
              <div className="right">
                <Skeleton className="w-24 h-10"/>
                <Skeleton className="price text-sm h-6 w-20"/>
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}