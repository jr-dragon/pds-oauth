import { createFileRoute } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";

export const Route = createFileRoute("/")({
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<main className="flex items-center justify-center min-h-screen">
			<Button type="button" size="lg" className="bg-bsky-blue hover:bg-bsky-blue-hover">
				以 Bluesky 帳號繼續
			</Button>
		</main>
	);
}
