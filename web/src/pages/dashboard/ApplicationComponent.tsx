import {WithStyles, withStyles} from "@mui/styles";
import {Application, Instance, RefBinding} from "../../axios";
import {Theme} from "@mui/material/styles";
import {Grid, TextField} from "@mui/material";
import InstanceComponent from "./InstanceComponent";

interface ApplicationProps extends WithStyles<typeof styles> {
    application: Application
    refBinding: RefBinding | undefined
    instances: Array<Instance>
}

const styles = (theme: Theme) => ({
    row: {
        padding: ".6em 1em",
        "&:hover": {
            backgroundColor: theme.palette.action.hover
        }
    },
    instances: {
        padding: ".6em 0"
    }
})

export default withStyles(styles)(function ApplicationComponent(props: ApplicationProps) {
    const {classes, application, refBinding, instances} = props
    return (
        <Grid className={classes.row} container>
            <Grid item xs={6} style={{display: "flex", alignItems: "flex-start"}}>
                <small>{application.name}</small>
            </Grid>
            <Grid item xs={6}>
                <TextField variant="standard" size="small" defaultValue={refBinding?.ref}/>
            </Grid>
            <Grid className={classes.instances} item xs={12}>
                {
                    instances.map(i => <InstanceComponent key={i.name} instance={i}/>)
                }
            </Grid>
        </Grid>
    )
})