import {Grid} from "@mui/material";
import {useEffect, useState} from "react";
import {Application, DefaultApiFp, Environment, RefBinding} from "../../axios";
import EnvironmentComponent from "./EnvironmentComponent";
import sum from "hash-sum";

const api = DefaultApiFp()

export default function DashboardPage() {
    const [environments, setEnvironments] = useState<Array<Environment>>([])
    const [applications, setApplications] = useState<Array<Application>>([])
    const [refBindings, setRefBindings] = useState<Array<RefBinding>>([])

    const [applicationsHash, setApplicationsHash] = useState<string>("")
    const [refBindingHashes, setRefBindingHashes] = useState<Map<string, string>>(new Map<string, string>())

    useEffect(() => {
        api.apiV1EnvironmentsGet()
            .then(request => request())
            .then(response => setEnvironments(response.data))

        api.apiV1ApplicationsGet()
            .then(request => request())
            .then(response => {
                setApplications(response.data)
                setApplicationsHash(sum(response.data))
            })

        api.apiV1RefBindingsGet()
            .then(request => request())
            .then(response => {
                setRefBindings(response.data)
                setRefBindingHashes(
                    Array.from(
                        response.data.reduce((set, item) => {
                            set.add(item.environment)
                            return set
                        }, new Set<string>())
                            .values()
                    )
                        .reduce((map, environment) => {
                            map.set(environment, sum(response.data.filter(r => r.environment === environment)))
                            return map
                        }, new Map<string, string>())
                )
            })

        return () => {
        }
    }, [])

    return (
        <Grid container spacing={2}>
            {
                environments.map(environment =>
                    <Grid key={`${environment.name}-${applicationsHash}-${refBindingHashes.get(environment.name)}`} item xs={12} sm={6} md={4} lg={3}>
                        <EnvironmentComponent
                            environment={environment}
                            applications={applications}
                            refBindings={refBindings.filter(r => r.environment === environment.name)}
                        />
                    </Grid>
                )
            }
        </Grid>
    )
}