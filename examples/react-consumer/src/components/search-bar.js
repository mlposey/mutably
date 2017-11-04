import React, { Component } from 'react';
import '../styles/search-bar.css';

/** An input form used for verb lookup */
class SearchBar extends Component {
    constructor(props) {
        super(props);
        this.state = { word: 'fiets' }
    }

    render() {
        return <input
            value={this.state.word}
            onChange={event => this.onSearchChange(event.target.value)} />;
    }

    /**
     * An event that is fired when text in the input form changes.
     * @param {string} word a word to look up; should not be empty or null
     */
    onSearchChange(word) {
        this.setState({ word });
        if (word.length === 0) return;        
        this.props.onSearchWord(word);
    }
}

export default SearchBar;