import {WithStyles, withStyles} from "@mui/styles";
import {Application, DefaultApiFp, Instance, InstancePod, RefBinding} from "../../axios";
import {Theme} from "@mui/material/styles";
import {CircularProgress, Grid, InputAdornment, TextField} from "@mui/material";
import InstanceComponent from "./InstanceComponent";
import {useSnackbar} from "notistack";

interface ApplicationProps extends WithStyles<typeof styles> {
    application: Application
    refBinding: RefBinding | undefined
    instances: Array<Instance>
    instancePods: Array<InstancePod>
    onRefBindingChanged: (refBinding: RefBinding) => void
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

const api = DefaultApiFp()

export default withStyles(styles)(function ApplicationComponent(props: ApplicationProps) {
    const {classes, application, refBinding, instances, instancePods, onRefBindingChanged} = props
    const {enqueueSnackbar, closeSnackbar} = useSnackbar();

    const onRefChanged = (newRef: string) => {
        if (refBinding?.ref !== newRef) {
            const newRefBinding = {...refBinding, ref: newRef} as RefBinding
            api.apiV1RefBindingsPost(newRefBinding)
                .then(request => request())
                .then(response => {
                    onRefBindingChanged(response.data)
                    enqueueSnackbar(`Ref ${newRef} has been deployed to ${refBinding?.environment} environment`, {variant: "success"})
                })
                .catch(() => {
                    enqueueSnackbar(`Ref ${newRef} has not been deployed to ${refBinding?.environment} environment`, {variant: "error"})
                })
        }
    }

    const deploying = refBinding?.ref && !instancePods.every(pod => pod.ref === refBinding.ref)
    return (
        <Grid className={classes.row} container>
            <Grid item xs={6} style={{display: "flex", alignItems: "flex-start"}}>
                <small>{application.name}</small>
            </Grid>
            <Grid item xs={6}>
                <TextField variant="standard" size="small"
                           defaultValue={refBinding?.ref}
                           onBlur={e => onRefChanged(e.target.value)}
                           InputProps={{
                               endAdornment: (deploying && <InputAdornment position="end">
                                   <CircularProgress size={16} />
                               </InputAdornment>)
                           }}/>
            </Grid>
            <Grid className={classes.instances} item xs={12}>
                {
                    instances.map(i => <InstanceComponent
                        key={i.name}
                        instance={i}
                        instancePods={instancePods.filter(instancePod => instancePod.parents?.includes(`${i.type}/${i.name}`) || false)}
                    />)
                }
            </Grid>
        </Grid>
    )
})