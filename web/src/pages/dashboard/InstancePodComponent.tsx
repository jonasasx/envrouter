import {Card, Chip, CircularProgress, Divider, Grid, List, ListItem} from "@mui/material";
import React, {useEffect, useState} from "react";
import {Theme} from "@mui/material/styles";
import {WithStyles, withStyles} from "@mui/styles";
import {Commit, DefaultApiFp, InstancePod} from "../../axios";
import {useSnackbar} from "notistack";

interface InstancePodProps extends WithStyles<typeof styles> {
    instancePod: InstancePod
}

const styles = (theme: Theme) => ({
    card: {
        marginBottom: "1em"
    }
})

const api = DefaultApiFp()

export default withStyles(styles)(function InstancePodComponent(props: InstancePodProps) {
    const {classes, instancePod} = props
    const [loading, setLoading] = useState(false)
    const [commit, setCommit] = useState<Commit | undefined>(undefined)
    const {enqueueSnackbar, closeSnackbar} = useSnackbar();
    useEffect(() => {
        if (instancePod.commitSha) {
            setLoading(true)
            api.apiV1GitApplicationsApplicationNameCommitsShaGet(instancePod.commitSha, instancePod.application)
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