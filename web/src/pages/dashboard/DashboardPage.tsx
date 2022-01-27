import {Grid} from "@mui/material";
import {useEffect, useState} from "react";
import {Application, DefaultApiFp, Environment, Instance, RefBinding} from "../../axios";
import EnvironmentComponent from "./EnvironmentComponent";

const api = DefaultApiFp()

export default function DashboardPage() {
    const [environments, setEnvironments] = useState<Array<Environment>>([])
    const [applications, setApplications] = useState<Array<Application>>([])
    const [refBindings, setRefBindings] = useState<Array<RefBinding>>([])
    const [instances, setInstances] = useState<Array<Instance>>([])

    useEffect(() => {

        Promise.all([
            api.apiV1EnvironmentsGet().then(request => request()),
            api.apiV1ApplicationsGet().then(request => request()),
            api.apiV1RefBindingsGet().then(request => request()),
            api.apiV1InstancesGet().then(request => request())
        ]).then(([envs, apps, refs, instances]) => {
            console.log({envs: envs.data, apps, refs});
            setEnvironments(envs.data);
            setApplications(apps.data);
            setRefBindings(refs.data);
            setInstances(instances.data);
        })

        return () => {
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
                        />
                    </Grid>
                )
            }
        </Grid>
    )
}