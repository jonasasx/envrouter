import {Card, Chip, CircularProgress, Divider, Grid, List, ListItem} from "@mui/material";
import React, {useEffect, useState} from "react";
import {Theme} from "@mui/material/styles";
import {WithStyles, withStyles} from "@mui/styles";
import {Commit, DefaultApiFp, Application, InstancePod, Ref} from "../../axios";
import {useSnackbar} from "notistack";

interface InstancePodProps extends WithStyles<typeof styles> {
    application: Application
    instancePod: InstancePod
    refsHeads: Array<Ref>
}

const styles = (theme: Theme) => ({
    card: {
        marginBottom: "1em"
    }
})

const api = DefaultApiFp()

export default withStyles(styles)(function InstancePodComponent(props: InstancePodProps) {
    const {classes, application, instancePod} = props
    const [loading, setLoading] = useState(false)
    const [commit, setCommit] = useState<Commit | undefined>(undefined)
    const {enqueueSnackbar, closeSnackbar} = useSnackbar();
    useEffect(() => {
        if (instancePod.commitSha && application.repositoryName) {
            setLoading(true)
            api.apiV1GitRepositoriesRepositoryNameCommitsShaGet(instancePod.commitSha, application.repositoryName!)
                .then(request => request())
                .then(response => setCommit(response.data))
                .then(() => setLoading(false))
                .catch(() => {
                    setLoading(false)
                    enqueueSnackbar(`Git fetching failed`, {variant: "error"})
                })
            return () => {
                closeSnackbar()
            }
        }
    }, [enqueueSnackbar, closeSnackbar])

    const rows: Array<{ key: string, value: any }> = [
        {key: "Pod name", value: instancePod.name},
        {key: "Application", value: instancePod.application},
        {key: "Environment", value: instancePod.environment},
        {key: "Shard", value: "s01"},
        {
            key: "Status",
            value: <Chip label={instancePod.phase} color={instancePod.phase === "Running" ? "success" : "warning"}
                         variant="outlined" size="small"/>
        },
        {key: "Created", value: instancePod.createdTime},
        {key: "Started", value: instancePod.startedTime},
        {key: "Branch", value: instancePod.ref || "-"},
        {key: "Commit", value: instancePod.commitSha},
        {key: "Author", value: loading ? <CircularProgress size={16}/> : commit?.author},
        {key: "Commit time", value: loading ? <CircularProgress size={16}/> : commit?.timestamp && new Date(Date.parse(commit?.timestamp)).toLocaleString()},
        {key: "Commit Message", value: loading ? <CircularProgress size={16}/> : commit?.message},
    ]
    return (
        <Card className={classes.card}>
            <List>
                {
                    rows.map((row, index) => (
                        <React.Fragment key={row.key}>
                            {index === 0 || <Divider/>}
                            <ListItem>
                                <Grid container>
                                    <Grid item xs={4}>{row.key}:</Grid>
                                    <Grid item xs={8} sx={{fontFamily: 'Monospace', fontSize: 12}}>{row.value}</Grid>
                                </Grid>
                            </ListItem>
                        </React.Fragment>
                    ))
                }
            </List>
        </Card>
    )
})