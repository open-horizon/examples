(function() {
    class Title extends React.Component {
        /*
        constructor(props) {
            super(props);

            this.state = {
                internalNum: props.num,
            };

            this.intervalId = null;
        }

        componentDidMount() {
            this.intervalId = setInterval(() => {
                this.setState((prevState, props) => ({
                    internalNum: prevState.internalNum + 1,
                }));
            }, 2000);
        }

        componentWillUnmount() {
            clearInterval(this.intervalId);
        }
        */

        render() {
            return React.createElement('h1', null, 'Hello Client!');
        }
    }
    const root = React.createElement(Title);
    ReactDOM.render(root, document.getElementById('react-root'));

    console.log('Hi there');
}());
