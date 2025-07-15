## Usage Doc according to requirements

<h2 id="note">Note</h2>
<ul>
<li><p>For all endpoints below you need to send the following header:
  <code>Authorization</code> : <code>Basic Base64encoded(apiusername:accessToken)</code></p>
<p>  <code>api_username and access token will be provided separately</code></p>
</li>
<li><p>Lists are created dynamically using a combination of language_preference + frequency.
So there is no need to manually create lists. When a new subscriber is created, three list is generated automatically, collective list, verified list and unverified list. </p>
</li>
</ul>
<h3 id="requirement-1-initialize-a-new-subscriber-list-on-aws-">Requirement#1: Initialize a new subscriber list on AWS.</h3>
<p>Solution:</p>
<ol>
<li>No need to create list specifically. Lists are automatically created when adding a subscriber according to the language parameter in the attributes</li>
<li>Created list is available in POST Subscriber API(as shown below)</li>
<li><p>To Fetch all lists the following end point can be used:</p>
<p> <b>GET /api/lists</b></p>
<pre><code class="lang-json"> Response:
 {
 <span class="hljs-string">"data"</span>: {
     <span class="hljs-string">"results"</span>: [
             {
                 <span class="hljs-string">"id"</span>: <span class="hljs-number">1</span>,
                 <span class="hljs-string">"created_at"</span>: <span class="hljs-string">"2020-02-10T23:07:16.194843+01:00"</span>,
                 <span class="hljs-string">"updated_at"</span>: <span class="hljs-string">"2020-03-06T22:32:01.118327+01:00"</span>,
                 <span class="hljs-string">"uuid"</span>: <span class="hljs-string">"ce13e971-c2ed-4069-bd0c-240e9a9f56f9"</span>,
                 <span class="hljs-string">"name"</span>: <span class="hljs-string">"Default list"</span>,
                 <span class="hljs-string">"type"</span>: <span class="hljs-string">"public"</span>,
                 <span class="hljs-string">"optin"</span>: <span class="hljs-string">"double"</span>,
                 <span class="hljs-string">"tags"</span>: [
                     <span class="hljs-string">"test"</span>
                 ],
                 <span class="hljs-string">"subscriber_count"</span>: <span class="hljs-number">2</span>
             },
             {
                 <span class="hljs-string">"id"</span>: <span class="hljs-number">2</span>,
                 <span class="hljs-string">"created_at"</span>: <span class="hljs-string">"2020-03-04T21:12:09.555013+01:00"</span>,
                 <span class="hljs-string">"updated_at"</span>: <span class="hljs-string">"2020-03-06T22:34:46.405031+01:00"</span>,
                 <span class="hljs-string">"uuid"</span>: <span class="hljs-string">"f20a2308-dfb5-4420-a56d-ecf0618a102d"</span>,
                 <span class="hljs-string">"name"</span>: <span class="hljs-string">"get"</span>,
                 <span class="hljs-string">"type"</span>: <span class="hljs-string">"private"</span>,
                 <span class="hljs-string">"optin"</span>: <span class="hljs-string">"single"</span>,
                 <span class="hljs-string">"tags"</span>: [],
                 <span class="hljs-string">"subscriber_count"</span>: <span class="hljs-number">0</span>
             }
         ],
     }
 }
</code></pre>
<hr/>

</li>
</ol>
<h3 id="requirement-2-enter-data-into-portal-database">Requirement#2: Enter Data into Portal Database</h3>
<p>Solution:</p>
<ol>
<li><p>Subscribers Data:
 <b>GET /api/subscribers</b></p>
<p> Response:</p>
<pre><code class="lang-json"> {
     <span class="hljs-attr">"data"</span>: {
         <span class="hljs-attr">"results"</span>: [
             {
                 <span class="hljs-attr">"id"</span>: <span class="hljs-number">1</span>,
                 <span class="hljs-attr">"created_at"</span>: <span class="hljs-string">"2020-02-10T23:07:16.199433+01:00"</span>,
                 <span class="hljs-attr">"updated_at"</span>: <span class="hljs-string">"2020-02-10T23:07:16.199433+01:00"</span>,
                 <span class="hljs-attr">"uuid"</span>: <span class="hljs-string">"ea06b2e7-4b08-4697-bcfc-2a5c6dde8f1c"</span>,
                 <span class="hljs-attr">"email"</span>: <span class="hljs-string">"john@example.com"</span>,
                 <span class="hljs-attr">"name"</span>: <span class="hljs-string">"John Doe"</span>,
                 <span class="hljs-attr">"attribs"</span>: {
                     <span class="hljs-attr">"city"</span>: <span class="hljs-string">"Bengaluru"</span>,
                     <span class="hljs-attr">"good"</span>: <span class="hljs-literal">true</span>,
                     <span class="hljs-attr">"type"</span>: <span class="hljs-string">"known"</span>
                 },
                 <span class="hljs-attr">"status"</span>: <span class="hljs-string">"enabled"</span>,
                 <span class="hljs-attr">"lists"</span>: [
                     {
                         <span class="hljs-attr">"subscription_status"</span>: <span class="hljs-string">"unconfirmed"</span>,
                         <span class="hljs-attr">"id"</span>: <span class="hljs-number">1</span>,
                         <span class="hljs-attr">"uuid"</span>: <span class="hljs-string">"ce13e971-c2ed-4069-bd0c-240e9a9f56f9"</span>,
                         <span class="hljs-attr">"name"</span>: <span class="hljs-string">"Default list"</span>,
                         <span class="hljs-attr">"type"</span>: <span class="hljs-string">"public"</span>,
                         <span class="hljs-attr">"tags"</span>: [
                             <span class="hljs-string">"test"</span>
                         ],
                         <span class="hljs-attr">"created_at"</span>: <span class="hljs-string">"2020-02-10T23:07:16.194843+01:00"</span>,
                         <span class="hljs-attr">"updated_at"</span>: <span class="hljs-string">"2020-02-10T23:07:16.194843+01:00"</span>
                     }
                 ]
             },
             {
                 <span class="hljs-attr">"id"</span>: <span class="hljs-number">2</span>,
                 <span class="hljs-attr">"created_at"</span>: <span class="hljs-string">"2020-02-18T21:10:17.218979+01:00"</span>,
                 <span class="hljs-attr">"updated_at"</span>: <span class="hljs-string">"2020-02-18T21:10:17.218979+01:00"</span>,
                 <span class="hljs-attr">"uuid"</span>: <span class="hljs-string">"ccf66172-f87f-4509-b7af-e8716f739860"</span>,
                 <span class="hljs-attr">"email"</span>: <span class="hljs-string">"quadri@example.com"</span>,
                 <span class="hljs-attr">"name"</span>: <span class="hljs-string">"quadri"</span>,
                 <span class="hljs-attr">"attribs"</span>: {
                     <span class="hljs-attr">"frequency_preferences"</span>: [
                         <span class="hljs-string">"3m"</span>
                     ],
                     <span class="hljs-attr">"language_preferences"</span>: [
                         <span class="hljs-string">"Hindi"</span>
                     ],
                     <span class="hljs-attr">"verified"</span>: <span class="hljs-literal">false</span>
                 },
                 <span class="hljs-attr">"status"</span>: <span class="hljs-string">"enabled"</span>,
                 <span class="hljs-attr">"lists"</span>: [
                     {
                         <span class="hljs-attr">"subscription_status"</span>: <span class="hljs-string">"unconfirmed"</span>,
                         <span class="hljs-attr">"id"</span>: <span class="hljs-number">1</span>,
                         <span class="hljs-attr">"uuid"</span>: <span class="hljs-string">"ce13e971-c2ed-4069-bd0c-240e9a9f56f9"</span>,
                         <span class="hljs-attr">"name"</span>: <span class="hljs-string">"Default list"</span>,
                         <span class="hljs-attr">"type"</span>: <span class="hljs-string">"public"</span>,
                         <span class="hljs-attr">"tags"</span>: [
                             <span class="hljs-string">"test"</span>
                         ],
                         <span class="hljs-attr">"created_at"</span>: <span class="hljs-string">"2020-02-10T23:07:16.194843+01:00"</span>,
                         <span class="hljs-attr">"updated_at"</span>: <span class="hljs-string">"2020-02-10T23:07:16.194843+01:00"</span>
                     }
                 ]
             },
             {
                 <span class="hljs-attr">"id"</span>: <span class="hljs-number">3</span>,
                 <span class="hljs-attr">"created_at"</span>: <span class="hljs-string">"2020-02-19T19:10:49.36636+01:00"</span>,
                 <span class="hljs-attr">"updated_at"</span>: <span class="hljs-string">"2020-02-19T19:10:49.36636+01:00"</span>,
                 <span class="hljs-attr">"uuid"</span>: <span class="hljs-string">"5d940585-3cc8-4add-b9c5-76efba3c6edd"</span>,
                 <span class="hljs-attr">"email"</span>: <span class="hljs-string">"sugar@example.com"</span>,
                 <span class="hljs-attr">"name"</span>: <span class="hljs-string">"sugar"</span>,
                 <span class="hljs-attr">"attribs"</span>: {
                     <span class="hljs-attr">"frequency_preferences"</span>: [
                         <span class="hljs-string">"3m"</span>
                     ],
                     <span class="hljs-attr">"language_preferences"</span>: [
                         <span class="hljs-string">"Hindi"</span>
                     ],
                     <span class="hljs-attr">"verified"</span>: <span class="hljs-literal">false</span>
                 },
                 <span class="hljs-attr">"status"</span>: <span class="hljs-string">"enabled"</span>,
                 <span class="hljs-attr">"lists"</span>: [
                 {
                         <span class="hljs-attr">"subscription_status"</span>: <span class="hljs-string">"confirmed"</span>,
                         <span class="hljs-attr">"subscription_created_at"</span>: <span class="hljs-string">"2025-06-11T00:22:19.461238+05:30"</span>,
                         <span class="hljs-attr">"subscription_updated_at"</span>: <span class="hljs-string">"2025-06-11T00:22:19.461238+05:30"</span>,
                         <span class="hljs-attr">"subscription_meta"</span>: {},
                         <span class="hljs-attr">"id"</span>: <span class="hljs-number">12</span>,
                         <span class="hljs-attr">"uuid"</span>: <span class="hljs-string">"4846c7fa-d0fa-4038-acf5-e3df4c4fe186"</span>,
                         <span class="hljs-attr">"name"</span>: <span class="hljs-string">"hindi_3m"</span>,
                         <span class="hljs-attr">"type"</span>: <span class="hljs-string">"public"</span>,
                         <span class="hljs-attr">"optin"</span>: <span class="hljs-string">"double"</span>,
                         <span class="hljs-attr">"tags"</span>: <span class="hljs-literal">null</span>,
                         <span class="hljs-attr">"description"</span>: <span class="hljs-string">"Auto-created shared list"</span>,
                         <span class="hljs-attr">"created_at"</span>: <span class="hljs-string">"2025-06-11T00:22:19.461238+05:30"</span>,
                         <span class="hljs-attr">"updated_at"</span>: <span class="hljs-string">"2025-06-11T00:22:19.461238+05:30"</span>
                       },
                       {
                         <span class="hljs-attr">"subscription_status"</span>: <span class="hljs-string">"confirmed"</span>,
                         <span class="hljs-attr">"subscription_created_at"</span>: <span class="hljs-string">"2025-06-11T00:22:19.461238+05:30"</span>,
                         <span class="hljs-attr">"subscription_updated_at"</span>: <span class="hljs-string">"2025-06-11T00:22:19.461238+05:30"</span>,
                         <span class="hljs-attr">"subscription_meta"</span>: {},
                         <span class="hljs-attr">"id"</span>: <span class="hljs-number">14</span>,
                         <span class="hljs-attr">"uuid"</span>: <span class="hljs-string">"782e86e3-2018-4bfa-a593-b0abae45fd0b"</span>,
                         <span class="hljs-attr">"name"</span>: <span class="hljs-string">"hindi_3m_unverified"</span>,
                         <span class="hljs-attr">"type"</span>: <span class="hljs-string">"public"</span>,
                         <span class="hljs-attr">"optin"</span>: <span class="hljs-string">"double"</span>,
                         <span class="hljs-attr">"tags"</span>: <span class="hljs-literal">null</span>,
                         <span class="hljs-attr">"description"</span>: <span class="hljs-string">"Auto-created unverified list"</span>,
                         <span class="hljs-attr">"created_at"</span>: <span class="hljs-string">"2025-06-11T00:22:19.461238+05:30"</span>,
                         <span class="hljs-attr">"updated_at"</span>: <span class="hljs-string">"2025-06-11T00:22:19.461238+05:30"</span>
                       }
                 ]
             }
         ],
         <span class="hljs-attr">"query"</span>: <span class="hljs-string">""</span>,
         <span class="hljs-attr">"total"</span>: <span class="hljs-number">3</span>,
         <span class="hljs-attr">"per_page"</span>: <span class="hljs-number">20</span>,
         <span class="hljs-attr">"page"</span>: <span class="hljs-number">1</span>
     }
 }
</code></pre>
</li>
<li><p>List Data
 <b>GET /api/lists</b></p>
<p> Response:</p>
<pre><code class="lang-json"> {
   <span class="hljs-attr">"data"</span>: {
     <span class="hljs-attr">"results"</span>: [
       {
         <span class="hljs-attr">"id"</span>: <span class="hljs-number">13</span>,
         <span class="hljs-attr">"created_at"</span>: <span class="hljs-string">"2025-06-11T00:22:19.461238+05:30"</span>,
         <span class="hljs-attr">"updated_at"</span>: <span class="hljs-string">"2025-06-11T00:22:19.461238+05:30"</span>,
         <span class="hljs-attr">"uuid"</span>: <span class="hljs-string">"4ae829db-a6ee-4d4d-914d-f2022f54998e"</span>,
         <span class="hljs-attr">"name"</span>: <span class="hljs-string">"hindi_3m_verified"</span>,
         <span class="hljs-attr">"type"</span>: <span class="hljs-string">"public"</span>,
         <span class="hljs-attr">"optin"</span>: <span class="hljs-string">"double"</span>,
         <span class="hljs-attr">"tags"</span>: [],
         <span class="hljs-attr">"description"</span>: <span class="hljs-string">"Auto-created verified list"</span>,
         <span class="hljs-attr">"subscriber_count"</span>: <span class="hljs-number">0</span>,
         <span class="hljs-attr">"subscriber_statuses"</span>: {},
         <span class="hljs-attr">"subscription_created_at"</span>: <span class="hljs-literal">null</span>,
         <span class="hljs-attr">"subscription_updated_at"</span>: <span class="hljs-literal">null</span>
       },
       {
         <span class="hljs-attr">"id"</span>: <span class="hljs-number">12</span>,
         <span class="hljs-attr">"created_at"</span>: <span class="hljs-string">"2025-06-11T00:22:19.461238+05:30"</span>,
         <span class="hljs-attr">"updated_at"</span>: <span class="hljs-string">"2025-06-11T00:22:19.461238+05:30"</span>,
         <span class="hljs-attr">"uuid"</span>: <span class="hljs-string">"4846c7fa-d0fa-4038-acf5-e3df4c4fe186"</span>,
         <span class="hljs-attr">"name"</span>: <span class="hljs-string">"hindi_3m"</span>,
         <span class="hljs-attr">"type"</span>: <span class="hljs-string">"public"</span>,
         <span class="hljs-attr">"optin"</span>: <span class="hljs-string">"double"</span>,
         <span class="hljs-attr">"tags"</span>: [],
         <span class="hljs-attr">"description"</span>: <span class="hljs-string">"Auto-created shared list"</span>,
         <span class="hljs-attr">"subscriber_count"</span>: <span class="hljs-number">1</span>,
         <span class="hljs-attr">"subscriber_statuses"</span>: {
           <span class="hljs-attr">"confirmed"</span>: <span class="hljs-number">1</span>
         },
         <span class="hljs-attr">"subscription_created_at"</span>: <span class="hljs-literal">null</span>,
         <span class="hljs-attr">"subscription_updated_at"</span>: <span class="hljs-literal">null</span>
       },
   }
 ]
</code></pre>
</li>
</ol>
<hr/>

<h3 id="requirement-3">Requirement#3</h3>
<ol>
<li><p>Create Subscriber and Subscribe to a list</p>
<p> 1.1 Create Subscriber.</p>
<p> Note: Lists are automatically created using language and frequency</p>
<p> Request:
 <b>POST /api/subscribers</b></p>
<pre><code class="lang-json"> {
    <span class="hljs-attr">"email"</span>:<span class="hljs-string">"subscriber@domain.com"</span>,
    <span class="hljs-attr">"name"</span>:<span class="hljs-string">"The Subscriber"</span>,
    <span class="hljs-attr">"status"</span>:<span class="hljs-string">"enabled"</span>,
    <span class="hljs-attr">"attribs"</span>:{
       <span class="hljs-attr">"frequency_preferences"</span>:[
          <span class="hljs-string">"3m"</span>
       ],
       <span class="hljs-attr">"language_preferences"</span>:[
          <span class="hljs-string">"Hindi"</span>
       ],
       <span class="hljs-attr">"verified"</span>:<span class="hljs-literal">false</span>
    }
 }
</code></pre>
<p> Response:</p>
<pre><code class="lang-json"> {
   <span class="hljs-attr">"data"</span>: {
     <span class="hljs-attr">"id"</span>: <span class="hljs-number">3</span>,
     <span class="hljs-attr">"created_at"</span>: <span class="hljs-string">"2019-07-03T12:17:29.735507+05:30"</span>,
     <span class="hljs-attr">"updated_at"</span>: <span class="hljs-string">"2019-07-03T12:17:29.735507+05:30"</span>,
     <span class="hljs-attr">"uuid"</span>: <span class="hljs-string">"eb420c55-4cfb-4972-92ba-c93c34ba475d"</span>,
     <span class="hljs-attr">"email"</span>: <span class="hljs-string">"subscriber@domain.com"</span>,
     <span class="hljs-attr">"name"</span>: <span class="hljs-string">"The Subscriber"</span>,
     <span class="hljs-attr">"attribs"</span>:{
       <span class="hljs-attr">"frequency_preferences"</span>:[
          <span class="hljs-string">"3m"</span>
       ],
       <span class="hljs-attr">"language_preferences"</span>:[
          <span class="hljs-string">"Hindi"</span>
       ],
       <span class="hljs-attr">"verified"</span>:<span class="hljs-literal">false</span>
     },
     <span class="hljs-attr">"status"</span>: <span class="hljs-string">"enabled"</span>,
     <span class="hljs-attr">"lists"</span>: [
         {
             <span class="hljs-attr">"subscription_status"</span>:<span class="hljs-string">"confirmed"</span>,
             <span class="hljs-attr">"subscription_created_at"</span>:<span class="hljs-string">"2025-07-15T21:31:37.250884+05:30"</span>,
             <span class="hljs-attr">"subscription_updated_at"</span>:<span class="hljs-string">"2025-07-15T21:31:37.250884+05:30"</span>,
             <span class="hljs-attr">"subscription_meta"</span>:{},
             <span class="hljs-attr">"id"</span>:<span class="hljs-number">12</span>,
             <span class="hljs-attr">"uuid"</span>:<span class="hljs-string">"4846c7fa-d0fa-4038-acf5-e3df4c4fe186"</span>,
             <span class="hljs-attr">"name"</span>:<span class="hljs-string">"hindi_3m"</span>,
             <span class="hljs-attr">"type"</span>:<span class="hljs-string">"public"</span>,
             <span class="hljs-attr">"optin"</span>:<span class="hljs-string">"double"</span>,
             <span class="hljs-attr">"tags"</span>:<span class="hljs-literal">null</span>,
             <span class="hljs-attr">"description"</span>:<span class="hljs-string">"Auto-created shared list"</span>,
             <span class="hljs-attr">"created_at"</span>:<span class="hljs-string">"2025-06-11T00:22:19.461238+05:30"</span>,
             <span class="hljs-attr">"updated_at"</span>:<span class="hljs-string">"2025-06-11T00:22:19.461238+05:30"</span>
         },
         {
             <span class="hljs-attr">"subscription_status"</span>:<span class="hljs-string">"confirmed"</span>,
             <span class="hljs-attr">"subscription_created_at"</span>:<span class="hljs-string">"2025-07-15T21:31:37.250884+05:30"</span>,
             <span class="hljs-attr">"subscription_updated_at"</span>:<span class="hljs-string">"2025-07-15T21:31:37.250884+05:30"</span>,
             <span class="hljs-attr">"subscription_meta"</span>:{},
             <span class="hljs-attr">"id"</span>:<span class="hljs-number">14</span>,
             <span class="hljs-attr">"uuid"</span>:<span class="hljs-string">"782e86e3-2018-4bfa-a593-b0abae45fd0b"</span>,
             <span class="hljs-attr">"name"</span>:<span class="hljs-string">"hindi_3m_unverified"</span>,
             <span class="hljs-attr">"type"</span>:<span class="hljs-string">"public"</span>,
             <span class="hljs-attr">"optin"</span>:<span class="hljs-string">"double"</span>,
             <span class="hljs-attr">"tags"</span>:<span class="hljs-literal">null</span>,
             <span class="hljs-attr">"description"</span>:<span class="hljs-string">"Auto-created unverified list"</span>,
             <span class="hljs-attr">"created_at"</span>:<span class="hljs-string">"2025-06-11T00:22:19.461238+05:30"</span>,
             <span class="hljs-attr">"updated_at"</span>:<span class="hljs-string">"2025-06-11T00:22:19.461238+05:30"</span>
         }
     ]
   }
 }
</code></pre>
</li>
</ol>
<hr/>

<h3 id="requirement-4-email-confirmation-update-subscriber-attributes">Requirement#4: Email Confirmation == Update subscriber attributes</h3>
<p>Request:
    <b>PUT PUT /api/subscribers/{subscriber_id}</b></p>
<pre><code>```json
{ 
   <span class="hljs-string">"attribs"</span>:{
      <span class="hljs-string">"frequency_preferences"</span>:[<span class="hljs-string">"3m"</span>],
      <span class="hljs-string">"language_preferences"</span>:[<span class="hljs-string">"English"</span>],  # example to show that changing langauge <span class="hljs-keyword">in</span> attribute also changes member ship
      <span class="hljs-string">"verified"</span>:true
   },
   <span class="hljs-string">"email"</span>: <span class="hljs-string">"subscriber9@domain.com"</span>
}
```
</code></pre><p>Response:</p>
<pre><code>```json
    {
      <span class="hljs-string">"data"</span>: {
        <span class="hljs-string">"id"</span>: <span class="hljs-number">3</span>,
        <span class="hljs-string">"created_at"</span>: <span class="hljs-string">"2019-07-03T12:17:29.735507+05:30"</span>,
        <span class="hljs-string">"updated_at"</span>: <span class="hljs-string">"2019-07-03T12:17:29.735507+05:30"</span>,
        <span class="hljs-string">"uuid"</span>: <span class="hljs-string">"eb420c55-4cfb-4972-92ba-c93c34ba475d"</span>,
        <span class="hljs-string">"email"</span>: <span class="hljs-string">"subscriber@domain.com"</span>,
        <span class="hljs-string">"name"</span>: <span class="hljs-string">"The Subscriber"</span>,
        <span class="hljs-string">"attribs"</span>:{
          <span class="hljs-string">"frequency_preferences"</span>:[
             <span class="hljs-string">"3m"</span>
          ],
          <span class="hljs-string">"language_preferences"</span>:[
             <span class="hljs-string">"English"</span>
          ],
          <span class="hljs-string">"verified"</span>:true
        },
        <span class="hljs-string">"status"</span>: <span class="hljs-string">"enabled"</span>,
        <span class="hljs-string">"lists"</span>: [
            {
                <span class="hljs-string">"subscription_status"</span>:<span class="hljs-string">"confirmed"</span>,
                <span class="hljs-string">"subscription_created_at"</span>:<span class="hljs-string">"2025-07-15T21:35:15.892767+05:30"</span>,
                <span class="hljs-string">"subscription_updated_at"</span>:<span class="hljs-string">"2025-07-15T21:35:15.892767+05:30"</span>,
                <span class="hljs-string">"subscription_meta"</span>:{},
                <span class="hljs-string">"id"</span>:<span class="hljs-number">9</span>,
                <span class="hljs-string">"uuid"</span>:<span class="hljs-string">"4f7b62ac-5fd4-4a1a-93d2-b06f86fc6741"</span>,
                <span class="hljs-string">"name"</span>:<span class="hljs-string">"english_3m"</span>,
                <span class="hljs-string">"type"</span>:<span class="hljs-string">"public"</span>,
                <span class="hljs-string">"optin"</span>:<span class="hljs-string">"double"</span>,
                <span class="hljs-string">"tags"</span>:null,
                <span class="hljs-string">"description"</span>:<span class="hljs-string">"Auto-created shared list"</span>,
                <span class="hljs-string">"created_at"</span>:<span class="hljs-string">"2025-06-11T00:22:02.221377+05:30"</span>,
                <span class="hljs-string">"updated_at"</span>:<span class="hljs-string">"2025-06-11T00:22:02.221377+05:30"</span>
            },
            {
                <span class="hljs-string">"subscription_status"</span>:<span class="hljs-string">"confirmed"</span>,
                <span class="hljs-string">"subscription_created_at"</span>:<span class="hljs-string">"2025-07-15T21:35:15.892767+05:30"</span>,
                <span class="hljs-string">"subscription_updated_at"</span>:<span class="hljs-string">"2025-07-15T21:35:15.892767+05:30"</span>,
                <span class="hljs-string">"subscription_meta"</span>:{},
                <span class="hljs-string">"id"</span>:<span class="hljs-number">10</span>,
                <span class="hljs-string">"uuid"</span>:<span class="hljs-string">"964f9e50-3247-40a0-8c6c-eae8067d5662"</span>,
                <span class="hljs-string">"name"</span>:<span class="hljs-string">"english_3m_verified"</span>,
                <span class="hljs-string">"type"</span>:<span class="hljs-string">"public"</span>,
                <span class="hljs-string">"optin"</span>:<span class="hljs-string">"double"</span>,
                <span class="hljs-string">"tags"</span>:null,
                <span class="hljs-string">"description"</span>:<span class="hljs-string">"Auto-created verified list"</span>,
                <span class="hljs-string">"created_at"</span>:<span class="hljs-string">"2025-06-11T00:22:02.221377+05:30"</span>,
                <span class="hljs-string">"updated_at"</span>:<span class="hljs-string">"2025-06-11T00:22:02.221377+05:30"</span>
            }
        ]
      }
    }
```
</code></pre><hr/>

<h3 id="requirement-4-aws-subscription-api-call">Requirement#4 : AWS Subscription API Call</h3>
<p>Solution: Already solved by requirement 2 solution</p>
<hr/>

<h3 id="requirement-5-manual-subscription-fallback">Requirement#5:  Manual Subscription Fallback</h3>
<p>Solution: Retry loop/Temprorary log of failed requests -&gt; Same POST /api/subscribers endpoint will be used</p>
<hr/>

<h3 id="requirement-6-unsubscribe-from-list">Requirement#6: Unsubscribe from List</h3>
<p>Solution:
    Request <b>PUT /api/subscribers/list</b></p>
<pre><code>```json
    {
       <span class="hljs-string">"ids"</span>:[
          <span class="hljs-number">1</span>,  # Subscriber Id
          <span class="hljs-number">2</span>,
          <span class="hljs-number">3</span>
       ],
       <span class="hljs-string">"action"</span>:<span class="hljs-string">"add"</span>, # add, remove, or unsubscribe
       <span class="hljs-string">"target_list_ids"</span>:[
          <span class="hljs-number">4</span>, # action will be performed on these list ids
          <span class="hljs-number">5</span>,
          <span class="hljs-number">6</span>
       ],
       <span class="hljs-string">"status"</span>:<span class="hljs-string">"confirmed"</span>
    }
```

Response: ```{data: true}```
</code></pre><hr/>

<h3 id="requirement-7-manual-unsubscribe-fallback">Requirement#7: Manual Unsubscribe Fallback</h3>
<p>Solution: Retry request on same endpoint, else log in fallback table</p>


<hr/>

<h3> Requirement#8: To send newsletter </h3>

Request:<br/>
    <b>POST &gt different-port &lt/proxy/send_campaign </b>
    <code>
        {
            "name":"name of newsletter",
            "subject": "subject of news letter",
            "type": "regualr",
            "content": "<html content goes here>"
            "lists": [listId1, listId2],
            "from_email": "noreply@bihariji.org",
            "from_name": "Shri Banke Bihariji"
        }
    </code>

<br/>Response: {success: true}<br/>
<code>
</code>

