(async () => {
	try {
		if (typeof mermaid === 'undefined') {
			return { svg: null, error: 'mermaid is not defined' };
		}

		mermaid.initialize({
			startOnLoad: false,
			securityLevel: 'loose',
			theme: 'default'
		});

		const id = 'mermaid-' + Math.random().toString(36).substr(2, 9);
		const result = await mermaid.render(id, {{.Code}});
		return { svg: result.svg || null, error: null };
	} catch (e) {
		return { svg: null, error: e.message || String(e) };
	}
})()
