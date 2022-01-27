import {Card, Chip, Divider, Grid, List, ListItem} from "@mui/material";
import React from "react";
import {Theme} from "@mui/material/styles";
import {WithStyles, withStyles} from "@mui/styles";
import {InstancePod} from "../../axios";

interface InstancePodProps extends WithStyles<typeof styles> {
    instancePod: InstancePod
}

const styles = (theme: Theme) => ({
    card: {
        marginBottom: "1em"
    }
})

export default withStyles(styles)(function InstancePodComponent(props: InstancePodProps) {
    const {classes, instancePod} = props

    const rows: Array<{ key: string, value: any }> = [
        {key: "Pod name", value: instancePod.name},
        {key: "Application", value: instancePod.application},
        {key: "Environment", value: instancePod.environment},
        {key: "Shard", value: "s01"},
        {key: "Status",
            value: <Chip label={instancePod.phase} color={instancePod.phase === "Running" ? "success" : "warning"}
                         variant="outlined" size="small"/>
        },
        {key: "Started", value: instancePod.startedTime},
        {key: "Branch", value: instancePod.ref || "-"},
        {key: "Commit", value: instancePod.commitSha},
        {key: "Author", value: "Ivan Volynkin"},
        {key: "Commit Message", value: "Init commit"},
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