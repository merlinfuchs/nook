import { ArrowRightIcon, TrendingUpIcon } from "lucide-react";
import { Button } from "../ui/button";
import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../ui/card";
import QuickstartDialog from "./QuickstartDialog";

export default function QuickstartCard() {
  return (
    <Card className="border-primary border-2 p-0 flex flex-col">
      <CardHeader className="mb-3 p-5 pb-0 flex-auto">
        <div className="flex items-center gap-5">
          <div className="size-12 border rounded-lg flex items-center justify-center flex-none">
            <TrendingUpIcon className="size-7" />
          </div>
          <div>
            <CardDescription className="mb-0.5 uppercase text-xs">
              Quick Start
            </CardDescription>
            <CardTitle className="text-xl">New around here?</CardTitle>
          </div>
        </div>
      </CardHeader>
      <CardFooter className="gap-3 flex-col md:flex-row p-5 pt-1">
        <QuickstartDialog>
          <Button className="w-full md:w-auto">
            <ArrowRightIcon className="size-4" />
            Get Started
          </Button>
        </QuickstartDialog>
      </CardFooter>
    </Card>
  );
}
