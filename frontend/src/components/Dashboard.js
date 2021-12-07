import AppBar from '@material-ui/core/AppBar';
import Box from '@material-ui/core/Box';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import React from "react";
import ApproveCardWithStyle from './Card.js';


class Dashboard extends React.Component {
    constructor(props) {
      super(props);
      this.state = {date: new Date()};
    }
  
    componentDidMount() {
      console.log('I was triggered during componentDidMount')
  
      fetch(window._env_.APPROVE_URL + '/requsts')
          .then(response => response.json())
          .then(data => this.setApprovelists(data));
    }
  
    setApprovelists(requests) {
      console.log(requests)

      let approveTemlist = []
      if(requests) {
        for(let i=0;i<requests.length;i++){
          let requestItem = {}

          requestItem.templateNS = requests[i].ApproveTemplate.split("/")[0]
          requestItem.templateName = requests[i].ApproveTemplate.split("/")[1]
          requestItem.requestNum = requests[i].Requests.length
          requestItem.requests = requests[i].Requests

          approveTemlist.push(requestItem)
        }
      }
      
      this.setState({
        date: new Date(),
        approveList: approveTemlist,
      })
    }

    render(){
      return (
        <Box sx={{ flexGrow: 1}}>
          <AppBar position="static">
            <Toolbar>
              <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                Tekton Manual Approve
              </Typography>
              <Button onClick={() => {
                  window.location.reload(false);
                }} color="inherit">Refresh</Button>
            </Toolbar>
          </AppBar>
          <Grid container spacing={2}>
            { this.state.approveList && this.state.approveList.map((item,index)=>{
              return <Grid item xs={3}><ApproveCardWithStyle approve={item}/></Grid>
            })}
          </Grid>
        </Box>
      )
    }
  }

export default Dashboard;
