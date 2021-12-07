import Avatar from '@material-ui/core/Avatar';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardHeader from '@material-ui/core/CardHeader';
import CardMedia from '@material-ui/core/CardMedia';
import { red } from '@material-ui/core/colors';
import IconButton from '@material-ui/core/IconButton';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/styles';
import PendingActionsIcon from '@mui/icons-material/PendingActions';
import Badge from '@mui/material/Badge';
import React from "react";
import { Link } from "react-router-dom";


const styles = theme => ({
  card: {
    maxWidth: 400,
  },
  media: {           
    height: 150,
    width: '300',
  },
});

class ApproveCard extends React.Component {
    constructor(props) {
      console.log(props)
      super(props);
      this.state = {ready: false};
    }

    componentDidMount() {
      console.log('I was triggered during componentDidMount in card components')
      console.log(this.props)
    }

    render(){
        return <Card sx={{ width: 330, margin: 2}}>
            <CardHeader
              avatar={
                <Avatar sx={{ bgcolor: red[500] }} aria-label="recipe">
                  RT
                </Avatar>
              }
              action={
                // <Link to={`/requests/${ this.props.approve }`}>
                <Link to={{  
                  pathname: "/requests/dummy", 
                  approve: this.props.approve }}>
                <IconButton aria-label="settings">
                  <Badge badgeContent={this.props.approve.requestNum} color="primary">
                      <PendingActionsIcon color="action" />
                  </Badge>
                </IconButton>
                </Link>
              }
              title="Approve Template"
              subheader="Request will attached on"
            />
            <CardMedia
              component="img"
              alt="approveTemplate"
              image="/images/header.png"
              className={this.props.classes.media}
            />
            <CardContent>
              <Typography variant="h6" gutterBottom>
              Approve template:
              </Typography>
              <Typography sx={{ mb: 1 }} color="text.secondary">
                Namespace: {this.props.approve.templateNS}
              </Typography>
              <Typography sx={{ mb: 1 }} color="text.secondary">
                Name: {this.props.approve.templateName}
              </Typography>
            </CardContent>
            {/* <CardActions>
              {!this.state.ready && (
                  <CircularProgress
                  size={30}
                  sx={{
                    color: green[500],
                    position: 'absolute',
                  }}
                />
              )

              }

              <Button href={window._env_.REPORT_URL + "reports/" + ".html"} variant="contained" disabled={false} size="small">CHECK REPORT</Button>
            </CardActions> */}
        </Card>
    }


}

const ApproveCardWithStyle = withStyles(styles)(ApproveCard);
export default ApproveCardWithStyle;