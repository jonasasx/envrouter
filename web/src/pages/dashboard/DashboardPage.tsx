import {Grid} from "@mui/material";
import {useEffect, useState} from "react";
import {Application, DefaultApiFp, Environment, Instance, InstancePod, RefBinding} from "../../axios";
import EnvironmentComponent from "./EnvironmentComponent";
import {PodEvent} from "../../sse/api";

const api = DefaultApiFp()

export default function DashboardPage() {
    const [environments, setEnvironments] = useState<Array<Environment>>([])
    const [applications, setApplications] = useState<Array<Application>>([])
    const [refBindings, setRefBindings] = useState<Array<RefBinding>>([])
    const [instances, setInstances] = useState<Array<Instance>>([])
    const [instancePods, setInstancePods] = useState<Array<InstancePod>>([])

    const onPodEvent: ((e: PodEvent) => void) = e => {
        const instancePod = e.item

        setInstancePods((currentInstancePods: Array<InstancePod>) => {
            const index = currentInstancePods.findIndex(i => i.name === instancePod.name)
            console.log('Event: ', e.event, '; Index: ', index, '; Name: ', instancePod.name, '; Phase: ', instancePod.phase)
            switch (e.event) {
                case "UPDATED":
                    if (index === -1) {
                        return [...currentInstancePods, instancePod]
                    }
                    return [
                        ...currentInstancePods.slice(0, index),
                        instancePod,
                        ...currentInstancePods.slice(index + 1)
                    ]
                case "DELETED":
                    return [
                        ...currentInstancePods.slice(0, index),
                        ...currentInstancePods.slice(index + 1)
                    ]
            }
        })
    }

    const updateRefBindingChanged = (newRefBinding: RefBinding) => {
        setRefBindings(currentRefBindings => {
            const index = refBindings.findIndex(r => r.environment === newRefBinding.environment && r.application === newRefBinding.application)
            return (index === -1) &&
                [...currentRefBindings, newRefBinding] ||
                [
                    ...currentRefBindings.slice(0, index),
                    newRefBinding,
                    ...currentRefBindings.slice(index + 1)
                ]
        })
    }

    useEffect(() => {
        const eventSource = new EventSource('http://localhost:8080/api/v1/subscription')
        eventSource.onmessage = e => onPodEvent(JSON.parse(e.data) as PodEvent)
        Promise.all([
            api.apiV1EnvironmentsGet().then(request => request()),
            api.apiV1ApplicationsGet().then(request => request()),
            api.apiV1RefBindingsGet().then(request => request()),
            api.apiV1InstancesGet().then(request => request()),
            api.apiV1InstancePodsGet().then(request => request())
        ]).then(([envs, apps, refs, instances, instancePods]) => {
            setRefBindings(refs.data);
            setEnvironments(envs.data.sort((a, b) => a.name.localeCompare(b.name)));
            setApplications(apps.data.sort((a, b) => a.name.localeCompare(b.name)));
            setInstances(instances.data);
            setInstancePods(instancePods.data.sort((a, b) => a.createdTime.localeCompare(b.createdTime)))
        })

        return () => {
            eventSource.close()
        }
    }, [])

    return (
        <Grid container spacing={2}>
            {
                environments.map(environment =>
                    <Grid
                        key={environment.name}
                        item xs={12} sm={6} md={4} lg={3}>
                        <EnvironmentComponent
                            environment={environment}
                            applications={applications}
                            refBindings={refBindings.filter(r => r.environment === environment.name)}
                            instances={instances.filter(i => i.environment === environment.name)}
                            instancePods={instancePods.filter(i => i.environment === environment.name)}
                            onRefBindingChanged={refBinding => updateRefBindingChanged(refBinding)}
                        />
                    </Grid>
                )
            }
        </Grid>
    )
}