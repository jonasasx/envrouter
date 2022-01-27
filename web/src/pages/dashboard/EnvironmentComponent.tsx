import {WithStyles, withStyles} from "@mui/styles";
import {Application, Environment, RefBinding} from "../../axios";
import {Theme} from "@mui/material/styles";
import {Card, CardContent, CardHeader} from "@mui/material";
import ApplicationComponent from "./ApplicationComponent";
import {Property} from "csstype"


interface EnvironmentProps extends WithStyles<typeof styles> {
    environment: Environment
    applications: Array<Application>
    refBindings: Array<RefBinding>
}

const styles = (theme: Theme) => ({
    cardHeader: {
        textAlign: 'left' as Property.TextAlign,
    },
    cardContent: {
        padding: 0,
        "&:last-child": {
            paddingBottom: 0
        }
    }
})

export default withStyles(styles)(function EnvironmentComponent(props: EnvironmentProps) {
    const {classes, environment, applications, refBindings} = props
    return (
        <Card>
            <CardHeader className={classes.cardHeader} title={environment.name}/>
            <CardContent className={classes.cardContent}>
                {
                    applications.map(application => <ApplicationComponent
                        application={application}
                        refBinding={refBindings.find(r => r.application === application.name)}
                    />)
                }
            </CardContent>
        </Card>
    )
})