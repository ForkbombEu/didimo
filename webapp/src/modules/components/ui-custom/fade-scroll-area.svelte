<!--
SPDX-FileCopyrightText: 2026 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import type { ClassValue } from 'svelte/elements';

	import { tick, type Snippet } from 'svelte';

	import { ScrollArea } from '@/components/ui/scroll-area';
	import { cn } from '@/components/ui/utils.js';

	type Props = {
		class?: ClassValue;
		contentClass?: ClassValue;
		/** Re-measure fades when this value changes (e.g. filtered list length). */
		refresh?: unknown;
		children?: Snippet;
	};

	let { class: className, contentClass, refresh, children }: Props = $props();

	let viewportRef = $state<HTMLElement | null>(null);
	let canScrollUp = $state(false);
	let canScrollDown = $state(false);

	function updateFades() {
		const viewport = viewportRef;
		if (!viewport) {
			canScrollUp = false;
			canScrollDown = false;
			return;
		}

		const { scrollTop, scrollHeight, clientHeight } = viewport;
		canScrollUp = scrollTop > 1;
		canScrollDown = scrollTop + clientHeight < scrollHeight - 1;
	}

	$effect(() => {
		const viewport = viewportRef;
		if (!viewport) return;

		updateFades();
		viewport.addEventListener('scroll', updateFades, { passive: true });

		const observer = new ResizeObserver(() => updateFades());
		observer.observe(viewport);
		for (const child of viewport.children) {
			observer.observe(child);
		}

		return () => {
			viewport.removeEventListener('scroll', updateFades);
			observer.disconnect();
		};
	});

	$effect(() => {
		void refresh;
		void tick().then(updateFades);
	});
</script>

<div class={cn('relative min-h-0 overflow-hidden', className)}>
	<ScrollArea
		bind:viewportRef
		type="always"
		class="size-full overflow-hidden"
		scrollbarYClasses="z-10 rounded-full bg-muted/80 [&_[data-slot=scroll-area-thumb]]:bg-muted-foreground/35"
	>
		<div class={cn('space-y-2 p-1 pr-5', contentClass)}>
			{@render children?.()}
		</div>
	</ScrollArea>

	{#if canScrollUp}
		<div
			class="pointer-events-none absolute start-0 end-3 top-0 z-[5] h-6 bg-gradient-to-b from-background to-transparent"
			aria-hidden="true"
		></div>
	{/if}

	{#if canScrollDown}
		<div
			class="pointer-events-none absolute start-0 end-3 bottom-0 z-[5] h-6 bg-gradient-to-t from-background to-transparent"
			aria-hidden="true"
		></div>
	{/if}
</div>
