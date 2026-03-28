import { createFileRoute } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { InputGroup, InputGroupAddon, InputGroupInput } from "@/components/ui/input-group";
import { useState } from "react";

export const Route = createFileRoute("/")({
	component: RouteComponent,
});

function RouteComponent() {
    const [username, setUsername] = useState("");

	return (
		<main className="flex items-center justify-center min-h-screen">
            <Dialog>
                <DialogTrigger asChild>
                    <Button type="button" size="lg" className="bg-bsky-blue hover:bg-bsky-blue-hover">
                        以 Bluesky 帳號繼續
                    </Button>
                </DialogTrigger>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>以 Bluesky 帳號繼續</DialogTitle>
                        <DialogDescription></DialogDescription>
                    </DialogHeader>
                    <form method="GET" action="/oauth/start">
                    <InputGroup>
                        <InputGroupAddon>@</InputGroupAddon>
                        <InputGroupInput placeholder="username.bsky.social" name="handle" value={username} onChange={e => setUsername(e.target.value)}/>
                    </InputGroup>
                    </form>
                </DialogContent>
            </Dialog>
		</main>
	);
}
