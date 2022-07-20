<script lang="ts">

	import { onMount,onDestroy } from 'svelte';
	enum keys{
		uuid,
		familyType,
		label,
		name,
		firmwareVersion
	}
	let entries: any[] = []
	let filter: Record<string,RegExp> = {}
	let interval;
	const compare = (a, b) => {
		const aFamily = a.params.device.familyType.toLowerCase()
		const aUUID = a.params.device.uuid.toLowerCase()
		const bFamily = b.params.device.familyType.toLowerCase()
		const bUUID = b.params.device.uuid.toLowerCase()
		const x = `${aFamily}-${aUUID}`
		const y = `${bFamily}-${bUUID}`
		if (x < y) {
		return -1
		}
		if (x > y) {
		return 1
		}
		return 0
	}

	const get =  () => {
      return fetch('/json').then((res)=>{
		if(res.status<300)return res.json()
		throw new Error("Could not reach server. Retrying in 1 second")	
	  }).then((data:Record<string,any>)=>{
		
		const newEntries = Object.entries(data)
      	.filter(([_ip,entry]) => {
		return Object.entries(filter).every(([key,value])=>{
			if(!value.test(entry.params.device[key]))return false
			return true
		})
      })
	  .map(([ip,entry])=> {
		entry.website=hasHTTP(entry["params"]["services"])?`http://${ip}`:""
		return entry
	  })
      .sort(compare)
      entries = [...newEntries] 

      })
    }

	const onChange = (key,event) => {
	   filter[key] = new RegExp(event.target.value)
    }
	
	onMount(async () => {
		interval= setInterval(async () => {get().catch((ex)=>
		console.log(ex))
	  }, 1000) 
		
	});
	onDestroy(() => clearInterval(interval));
	const hasHTTP = services => {
      return services && services.filter(s => s.type === 'http').length
    }
	

</script>

<main>
	<div>
		<h2>Filter</h2>
		{#each Object.values(keys).filter((el)=>typeof(el)==="string") as key,i}
			<div class="filterDiv">
				<label for={`filter${i}`} style="width:400px">{key}</label>
				<input id={`filter${i}`} type="text" name="filter" on:change={(ev)=>onChange(key,ev)} />
			</div>
		{/each}
		<table>
		  <thead>
			<tr>
			  <th />
				{#each Object.values(keys).filter((el)=>typeof(el)==="string") as key}
					<th>{key}</th>
				{/each}
			  <th>Website</th>
			</tr>
		  </thead>
		  <tbody>
			{#each entries as entry,i}
				<tr >
					<td class="index">#{i + 1}</td>
					{#each Object.values(keys).filter((el)=>typeof(el)==="string") as key}
						<td>{entry.params.device[key]}</td>
					{/each}
					<td><a href={entry.website}>{entry.website}</a></td>
				</tr>
			{/each}
		  </tbody>
		</table>
	  </div>
</main>

<style>
.filterDiv{
	width: 300px;
    display: flex;
    justify-content: space-between;
}
table {
  border-collapse: collapse;
}

th {
  text-align: left;
  padding: 0.35rem 0.75rem;
}

td {
  padding: 0.35rem 0.75rem;
}

tbody tr:hover {
  background-color: rgba(0, 0, 0, .075);
}

.index {
  color: #6f6f6f;
}
</style>