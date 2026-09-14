<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import type { Record } from '$lib/pipeline/device';

	import { browser } from '$app/environment';
	import { Pipeline } from '$lib';

	import type { PipelinesResponse } from '@/pocketbase/types';

	import Dialog from '@/components/ui-custom/dialog.svelte';

	import DeviceSelectInput from './device-select-input.svelte';

	//

	type Props = {
		open?: boolean;
		pipeline: PipelinesResponse;
		title?: string;
		description?: string;
		onSelect?: (device: Record) => void;
	};

	let {
		open = $bindable(false),
		pipeline,
		title = 'Select device',
		description = 'Select a device to execute the pipeline',
		onSelect
	}: Props = $props();

	//

	function handleSelect(device: Record) {
		Pipeline.Device.Binding.set(pipeline, device);
		open = false;
		onSelect?.(device);
	}

	//

	let currentDevicePath = $derived.by(() => {
		if (!browser) return undefined;
		return Pipeline.Device.Binding.get(pipeline.id);
	});

	$effect(() => {
		if (!open) return;
		void Pipeline.Device.Catalog.refresh();
	});
</script>

<Dialog
	bind:open
	{title}
	{description}
	hideTrigger
	contentClass="max-h-[calc(100dvh-2rem)] overflow-hidden"
>
	{#snippet content()}
		<DeviceSelectInput
			presentation="run"
			onSelect={handleSelect}
			selectedDevice={currentDevicePath}
			scrollable
		/>
	{/snippet}
</Dialog>
